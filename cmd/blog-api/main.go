package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cicd2jenkins/internal/app"
	"cicd2jenkins/internal/config"
)

func main() {
	cfg := config.Load()

	server, cleanup, err := app.NewServer(cfg)
	if err != nil {
		log.Fatalf("bootstrap application: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("%s listening on %s", cfg.AppName, cfg.Server.BindAddress())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server stopped unexpectedly: %v", err)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	if cleanup != nil {
		if err := cleanup(shutdownCtx); err != nil {
			log.Printf("cleanup failed: %v", err)
		}
	}
}
