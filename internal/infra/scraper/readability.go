package scraper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	readability "codeberg.org/readeck/go-readability/v2"
	"github.com/microcosm-cc/bluemonday"
)

// envAllowLoopback, when set to "1", lets the dialer reach loopback / private
// destinations. Test-only escape hatch — enabled by cmd/testserver and by the
// scraper's own httptest-based unit tests. Production callers MUST NOT set it.
const envAllowLoopback = "LITEREADER_SCRAPE_ALLOW_LOOPBACK"

const (
	defaultTimeout    = 10 * time.Second
	maxBodyBytes      = 5 * 1024 * 1024
	maxRedirects      = 5
	maxConcurrent     = 4
	userAgent         = "lite-reader/1.0 (+https://github.com/cubny/lite-reader)"
	dialTimeout       = 5 * time.Second
	keepAliveInterval = 30 * time.Second
)

// allowedSchemes is the explicit scheme allowlist. Items in feeds occasionally
// have file://, gopher://, javascript:, etc. in the link field — refuse them.
var allowedSchemes = map[string]struct{}{
	"http":  {},
	"https": {},
}

// allowedMIMETypes is the explicit content-type allowlist for response bodies.
// Used after parsing the response Content-Type via mime.ParseMediaType.
var allowedMIMETypes = map[string]struct{}{
	"text/html":             {},
	"application/xhtml+xml": {},
}

// errBlockedDestination is returned when the resolved IP for a fetch target
// is in a range we refuse to dial (loopback, RFC1918, link-local, etc.).
var errBlockedDestination = errors.New("scrape: blocked destination")

// Readability fetches an article URL and returns the sanitized main-content
// HTML. It implements item.Scraper. SSRF is mitigated by:
//   - URL scheme allowlist (http, https only)
//   - DialContext that resolves the host and rejects loopback / private /
//     link-local / multicast / unspecified IPs (catches the initial fetch
//     and every redirect, since Go opens a fresh connection per host)
//   - Bounded redirect chain
type Readability struct {
	client    *http.Client
	sanitizer *bluemonday.Policy
	sem       chan struct{}
}

func NewReadability() *Readability {
	dial := safeDialContext
	if os.Getenv(envAllowLoopback) == "1" {
		dial = (&net.Dialer{Timeout: dialTimeout, KeepAlive: keepAliveInterval}).DialContext
	}
	transport := &http.Transport{
		DialContext:           dial,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          16,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{
		Timeout:   defaultTimeout,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return fmt.Errorf("scrape: too many redirects (>%d)", maxRedirects)
			}
			if _, ok := allowedSchemes[req.URL.Scheme]; !ok {
				return fmt.Errorf("scrape: redirect to disallowed scheme %q", req.URL.Scheme)
			}
			return nil
		},
	}
	return &Readability{
		client:    client,
		sanitizer: bluemonday.UGCPolicy(),
		sem:       make(chan struct{}, maxConcurrent),
	}
}

// Scrape downloads rawURL, extracts the main article body, and returns the
// sanitized HTML. Errors are returned as-is so callers can persist a non-ok
// scrape_status without retrying in a tight loop.
func (s *Readability) Scrape(ctx context.Context, rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return "", fmt.Errorf("scrape: invalid url %q", rawURL)
	}
	if _, ok := allowedSchemes[parsed.Scheme]; !ok {
		return "", fmt.Errorf("scrape: disallowed scheme %q", parsed.Scheme)
	}

	// Bound concurrent scrapes so a logged-in user can't fan out parallel
	// 5MB downloads. Honor ctx cancellation while waiting for a slot.
	select {
	case s.sem <- struct{}{}:
		defer func() { <-s.sem }()
	case <-ctx.Done():
		return "", ctx.Err()
	}

	raw, err := s.fetchBody(ctx, rawURL)
	if err != nil {
		return "", err
	}
	return s.extract(raw, parsed)
}

func (s *Readability) fetchBody(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("scrape: build request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("scrape: fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("scrape: unexpected status %d", resp.StatusCode)
	}
	if ctErr := checkContentType(resp.Header.Get("Content-Type")); ctErr != nil {
		return nil, ctErr
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes+1))
	if err != nil {
		return nil, fmt.Errorf("scrape: read body: %w", err)
	}
	if len(raw) > maxBodyBytes {
		return nil, errors.New("scrape: body too large")
	}
	return raw, nil
}

func checkContentType(ct string) error {
	if ct == "" {
		return nil
	}
	mediaType, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return fmt.Errorf("scrape: invalid content-type %q: %w", ct, err)
	}
	if _, ok := allowedMIMETypes[mediaType]; !ok {
		return fmt.Errorf("scrape: disallowed content-type %q", mediaType)
	}
	return nil
}

func (s *Readability) extract(raw []byte, pageURL *url.URL) (string, error) {
	article, err := readability.FromReader(bytes.NewReader(raw), pageURL)
	if err != nil {
		return "", fmt.Errorf("scrape: extract: %w", err)
	}
	var rendered bytes.Buffer
	if err := article.RenderHTML(&rendered); err != nil {
		return "", fmt.Errorf("scrape: render: %w", err)
	}
	if strings.TrimSpace(rendered.String()) == "" {
		return "", errors.New("scrape: empty article content")
	}
	return s.sanitizer.Sanitize(rendered.String()), nil
}

// safeDialContext resolves the host and refuses to dial any IP in a
// non-routable / private / loopback / link-local / multicast range. This is
// the SSRF chokepoint and runs for both the initial request and every
// redirect (Go opens a fresh connection per host).
func safeDialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
	if err != nil {
		return nil, err
	}
	dialer := &net.Dialer{Timeout: dialTimeout, KeepAlive: keepAliveInterval}
	var lastErr error
	for _, ip := range ips {
		if isBlockedIP(ip) {
			lastErr = fmt.Errorf("%w: %s", errBlockedDestination, ip)
			continue
		}
		conn, dialErr := dialer.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
		if dialErr == nil {
			return conn, nil
		}
		lastErr = dialErr
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("scrape: no usable address for %s", host)
	}
	return nil, lastErr
}

func isBlockedIP(ip net.IP) bool {
	return ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsInterfaceLocalMulticast() ||
		ip.IsMulticast() ||
		ip.IsUnspecified()
}
