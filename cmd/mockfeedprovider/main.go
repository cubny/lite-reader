package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal/testserver"
)

const defaultFeedProviderPort = 3002

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	port := defaultFeedProviderPort
	if raw := os.Getenv("MOCK_FEED_PROVIDER_PORT"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			port = parsed
		} else {
			log.Warnf("invalid MOCK_FEED_PROVIDER_PORT %q, falling back to %d", raw, defaultFeedProviderPort)
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mockServer := testserver.NewMockFeedServer(port)
	if err := mockServer.Start(); err != nil {
		log.Fatalf("failed to start mock feed provider: %v", err)
	}

	log.Infof("Mock feed provider running on port %d", port)

	<-ctx.Done()

	if err := mockServer.Stop(); err != nil {
		log.Errorf("failed to stop mock feed provider cleanly: %v", err)
	}
}
