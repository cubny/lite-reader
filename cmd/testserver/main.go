package main

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal"
	"github.com/cubny/lite-reader/internal/testserver"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	ctx, cancel := context.WithCancel(context.Background())

	// Start mock feed server on port 3001
	mockServer := testserver.NewMockFeedServer(3001)
	if err := mockServer.Start(); err != nil {
		log.Fatalf("failed to start mock feed server: %v", err)
	}
	
	// Give the mock server time to start
	time.Sleep(500 * time.Millisecond)
	log.Info("Mock feed server started on port 3001")

	// Start main application
	runMigration := true
	app, err := internal.Init(ctx, runMigration)
	if err != nil {
		log.Fatalf("failed to initiate App, %v", err)
	}

	internal.WaitTermination()
	cancel()

	// Stop mock feed server
	if err = mockServer.Stop(); err != nil {
		log.Errorf("failed to stop mock feed server: %v", err)
	}

	// Stop main application
	if err = app.Stop(); err != nil {
		log.Errorf("failed to stop the app gracefully, %v", err)
	}
}
