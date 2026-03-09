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

	"tarotdaily/backend/internal/config"
	mongoplatform "tarotdaily/backend/internal/platform/mongo"
	httptransport "tarotdaily/backend/internal/transport/http"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mongoCtx, mongoCancel := context.WithTimeout(rootCtx, cfg.MongoConnectTimeout)
	defer mongoCancel()

	mongoClient, err := mongoplatform.Connect(mongoCtx, cfg)
	if err != nil {
		log.Fatalf("mongo connect failed: %v", err)
	}
	defer func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(disconnectCtx); err != nil {
			log.Printf("mongo disconnect warning: %v", err)
		}
	}()

	router := httptransport.NewRouter(cfg, mongoClient)

	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTPReadHeaderTimeout,
		ReadTimeout:       cfg.HTTPReadTimeout,
		WriteTimeout:      cfg.HTTPWriteTimeout,
		IdleTimeout:       cfg.HTTPIdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("api server started on :%s (env=%s)", cfg.AppPort, cfg.AppEnv)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-rootCtx.Done():
		log.Printf("shutdown signal received")
	case err := <-errCh:
		log.Fatalf("server crashed: %v", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown warning: %v", err)
	}

	log.Printf("api server stopped")
}
