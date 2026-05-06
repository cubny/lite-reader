package testserver

import (
	"embed"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//go:embed fixtures/*.xml
var feedFixtures embed.FS

// MockFeedServer serves mock RSS and Atom feeds for testing
type MockFeedServer struct {
	server *http.Server
	port   int
}

// NewMockFeedServer creates a new mock feed server
func NewMockFeedServer(port int) *MockFeedServer {
	return &MockFeedServer{
		port: port,
	}
}

// Start starts the mock feed server
func (m *MockFeedServer) Start() error {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Serve feed fixtures
	mux.HandleFunc("/feeds/", m.serveFeed)

	m.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", m.port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Infof("Mock feed server starting on port %d", m.port)
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("Mock feed server error: %v", err)
		}
	}()

	return nil
}

// Stop stops the mock feed server
func (m *MockFeedServer) Stop() error {
	if m.server != nil {
		return m.server.Close()
	}
	return nil
}

// GetURL returns the base URL for a feed
func (m *MockFeedServer) GetURL(feedName string) string {
	return fmt.Sprintf("http://localhost:%d/feeds/%s", m.port, feedName)
}

func (m *MockFeedServer) serveFeed(w http.ResponseWriter, r *http.Request) {
	// Extract feed name from path
	feedName := strings.TrimPrefix(r.URL.Path, "/feeds/")
	if feedName == "" {
		http.Error(w, "Feed name required", http.StatusBadRequest)
		return
	}

	// Ensure .xml extension
	if !strings.HasSuffix(feedName, ".xml") {
		feedName += ".xml"
	}

	// Read feed from fixtures
	feedPath := filepath.Join("fixtures", feedName)
	content, err := feedFixtures.ReadFile(feedPath)
	if err != nil {
		log.Warnf("Feed not found: %s", feedName)
		http.Error(w, "Feed not found", http.StatusNotFound)
		return
	}

	// Set appropriate content type
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	// content is a fixed RSS/Atom XML fixture loaded from the testserver fixtures directory.
	_, _ = w.Write(content) //nolint:gosec // G705: test-only mock feed server
}
