package main

import (
	"context"
	"github/devAshu12/learning_platform_GO_backend/pkg/config"
	"github/devAshu12/learning_platform_GO_backend/pkg/db"
	"github/devAshu12/learning_platform_GO_backend/pkg/handlers"
	"github/devAshu12/learning_platform_GO_backend/pkg/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	config.LoadConfig()
	db.InitDB()
	s := server.NewServer()

	// Initialize queue and start queue processor
	handlers.Quit = make(chan struct{})
	go handlers.ProcessQueue(10*time.Second, 10) // 10-second interval or 5 items to batch

	// Channel to listen for interrupt or termination signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGABRT)

	go func() {
		<-shutdown
		log.Println("Shutting down server...")

		// Signal ProcessQueue to stop
		close(handlers.Quit)

		// Wait briefly to allow for the final queue flush before shutdown
		time.Sleep(5 * time.Second)

		// Gracefully shut down the HTTP server
		if err := s.Shutdown(context.Background()); err != nil {
			log.Fatalf("Server forced to shutdown: %s", err)
		}
	}()

	log.Println("Starting server on :8080")
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
