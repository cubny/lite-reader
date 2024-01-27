package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	ctx, cancel := context.WithCancel(context.Background())

	app, err := internal.Init(ctx)
	if err != nil {
		log.Fatalf("failed to initiate App, %v", err)
	}

	internal.WaitTermination()
	cancel()

	if err = app.Stop(); err != nil {
		log.Errorf("failed to stop the app gracefully, %v", err)
	}
}
