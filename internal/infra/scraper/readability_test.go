package scraper

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	s := NewReadability()
	_, err := s.Scrape(context.Background(), srv.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status 404")
}

func TestReadability_NonHTML(t *testing.T) {
	srv := newTestServer(t, http.StatusOK, "application/pdf", "%PDF-1.4 ...")
	defer srv.Close()

	s := NewReadability()
	_, err := s.Scrape(context.Background(), srv.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "non-html")
}

func TestReadability_InvalidURL(t *testing.T) {
	s := NewReadability()
	_, err := s.Scrape(context.Background(), "not-a-url")
	require.Error(t, err)
}

func TestReadability_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_, _ = w.Write([]byte(sampleArticle))
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	s := NewReadability()
	_, err := s.Scrape(ctx, srv.URL)
	require.Error(t, err)
}
