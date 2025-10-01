package main

import (
	"context"
	"hexagonal/internal/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Bootstrap the application context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Initialize the bootstrap struct
	appInstance, err := app.Initialize(ctx)
	if err != nil {
		log.Fatalf("Error during initialization: %v", err)
	}

	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go appInstance.GracefulShutdown(signals, done)

	appInstance.Start(ctx)
	<-done

	cancel()

	log.Println("Waiting for background jobs to finish their works...")
	appInstance.Wait()

	log.Println("Server exited gracefully")
}
