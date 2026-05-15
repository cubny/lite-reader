package scraper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	readability "codeberg.org/readeck/go-readability/v2"
	"github.com/microcosm-cc/bluemonday"
)

const (
	defaultTimeout = 10 * time.Second
	maxBodyBytes   = 5 * 1024 * 1024
	userAgent      = "lite-reader/1.0 (+https://github.com/cubny/lite-reader)"
)

// Readability fetches an article URL and returns the sanitized main-content
// HTML. It implements item.Scraper.
type Readability struct {
	client    *http.Client
	sanitizer *bluemonday.Policy
}

func NewReadability() *Readability {
	return &Readability{
		client:    &http.Client{Timeout: defaultTimeout},
		sanitizer: bluemonday.UGCPolicy(),
	}
}

// Scrape downloads rawURL, extracts the main article body, and returns the
// sanitized HTML. Errors are returned as-is so callers can persist a non-ok
// scrape_status without retrying in a tight loop.
func (s *Readability) Scrape(ctx context.Context, rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("scrape: invalid url %q", rawURL)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("scrape: build request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("scrape: fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("scrape: unexpected status %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "" && !strings.Contains(ct, "html") {
		return "", fmt.Errorf("scrape: non-html content type %q", ct)
	}

	body := io.LimitReader(resp.Body, maxBodyBytes+1)
	raw, err := io.ReadAll(body)
	if err != nil {
		return "", fmt.Errorf("scrape: read body: %w", err)
	}
	if len(raw) > maxBodyBytes {
		return "", errors.New("scrape: body too large")
	}

	article, err := readability.FromReader(bytes.NewReader(raw), parsed)
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
