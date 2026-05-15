package scraper

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// allowLoopback enables the test-only escape hatch so httptest servers
// (always loopback) are reachable. Reset by t.Setenv after the test.
func allowLoopback(t *testing.T) {
	t.Helper()
	t.Setenv(envAllowLoopback, "1")
}

const sampleArticle = `<!doctype html><html><head><title>Sample</title></head><body>
<header><nav>Top nav junk</nav></header>
<article>
  <h1>How to brew coffee</h1>
  <p>Coffee is a beloved beverage enjoyed worldwide. This short guide walks through the essentials of brewing
  a great cup at home, from grind size to water temperature, with notes on common pitfalls and how to avoid them.</p>
  <p>Start with fresh beans, grind just before brewing, and measure your dose. A 1:16 ratio of coffee to water
  is a sensible starting point. Adjust to taste from there as you dial in your preferred extraction strength.</p>
  <p>Water temperature matters more than people think. Aim for somewhere between 90 and 96 degrees Celsius. Cooler
  water under-extracts and tastes sour; water that is too hot scorches the grounds and produces bitterness.</p>
  <script>alert('xss')</script>
</article>
<footer>Footer junk</footer>
</body></html>`

func newTestServer(t *testing.T, status int, contentType, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func TestReadability_Success(t *testing.T) {
	srv := newTestServer(t, http.StatusOK, "text/html; charset=utf-8", sampleArticle)
	defer srv.Close()

	allowLoopback(t)
	s := NewReadability()
	out, err := s.Scrape(context.Background(), srv.URL)
	require.NoError(t, err)
	assert.Contains(t, out, "How to brew coffee")
	assert.NotContains(t, out, "<script", "scripts must be stripped")
	assert.NotContains(t, strings.ToLower(out), "alert(", "script bodies must be stripped")
}

func TestReadability_Non2xx(t *testing.T) {
	srv := newTestServer(t, http.StatusNotFound, "text/html", "not found")
	defer srv.Close()

	allowLoopback(t)
	s := NewReadability()
	_, err := s.Scrape(context.Background(), srv.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status 404")
}

func TestReadability_NonHTML(t *testing.T) {
	srv := newTestServer(t, http.StatusOK, "application/pdf", "%PDF-1.4 ...")
	defer srv.Close()

	allowLoopback(t)
	s := NewReadability()
	_, err := s.Scrape(context.Background(), srv.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "disallowed content-type")
}

func TestReadability_InvalidURL(t *testing.T) {
	allowLoopback(t)
	s := NewReadability()
	_, err := s.Scrape(context.Background(), "not-a-url")
	require.Error(t, err)
}

func TestReadability_DisallowedScheme(t *testing.T) {
	// No allowLoopback — scheme check happens before any network use.
	s := NewReadability()
	for _, scheme := range []string{"file", "ftp", "gopher", "javascript"} {
		t.Run(scheme, func(t *testing.T) {
			_, err := s.Scrape(context.Background(), scheme+"://example.com/x")
			require.Error(t, err)
			assert.Contains(t, err.Error(), "disallowed scheme")
		})
	}
}

func TestReadability_BlocksLoopback(t *testing.T) {
	// Default config (no env override) — httptest binds to loopback, so the
	// safe dialer must refuse it.
	srv := newTestServer(t, http.StatusOK, "text/html", sampleArticle)
	defer srv.Close()

	s := NewReadability()
	_, err := s.Scrape(context.Background(), srv.URL)
	require.Error(t, err)
	assert.True(t,
		errors.Is(err, errBlockedDestination) || strings.Contains(err.Error(), "blocked destination"),
		"expected blocked-destination error, got %v", err)
}

func TestReadability_BlocksRedirectToLoopback(t *testing.T) {
	// Public-looking host that redirects to loopback should be refused.
	// We exercise this via a loopback server that 302s to itself: when the
	// dialer is the safe one, both legs are refused. A "real" public→private
	// redirect can't be exercised without external DNS, but the redirect
	// path runs through the same DialContext, so the guard is the same code.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://127.0.0.1:1/", http.StatusFound)
	}))
	defer srv.Close()

	s := NewReadability()
	_, err := s.Scrape(context.Background(), srv.URL)
	require.Error(t, err)
}

func TestIsBlockedIP(t *testing.T) {
	cases := map[string]bool{
		"127.0.0.1":     true,
		"::1":           true,
		"10.0.0.1":      true,
		"172.16.0.1":    true,
		"192.168.0.1":   true,
		"169.254.0.1":   true, // link-local (cloud metadata)
		"0.0.0.0":       true,
		"224.0.0.1":     true, // multicast
		"8.8.8.8":       false,
		"1.1.1.1":       false,
		"93.184.216.34": false, // example.com
	}
	for s, want := range cases {
		t.Run(s, func(t *testing.T) {
			assert.Equal(t, want, isBlockedIP(net.ParseIP(s)))
		})
	}
}

func TestReadability_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_, _ = w.Write([]byte(sampleArticle))
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	allowLoopback(t)
	s := NewReadability()
	_, err := s.Scrape(ctx, srv.URL)
	require.Error(t, err)
}
