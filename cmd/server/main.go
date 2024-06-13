package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wojciechpawlinow/eth-parser/internal/infrastructure/container"
	"github.com/wojciechpawlinow/eth-parser/internal/infrastructure/httpserver"
)

func main() {
	// build dependencies
	ctn := container.New()

	errChan := make(chan error, 1)

	// create and run HTTP server
	s := httpserver.Run(ctn, errChan)

	// wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds
	// use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("initiating graceful shutdown...")
	case err := <-errChan:
		log.Println("server error: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// gracefully shut down the server
	if err := s.Shutdown(ctx); err != nil {
		log.Println("server shutdown failed: ", err)
	}
}
