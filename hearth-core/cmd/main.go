package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mahaveer86619/Hearth/pkg/config"
	"github.com/Mahaveer86619/Hearth/pkg/db"
	"github.com/Mahaveer86619/Hearth/pkg/logger"
	"github.com/Mahaveer86619/Hearth/pkg/services"
	"github.com/Mahaveer86619/Hearth/pkg/services/ingestion_pipeline"
	"github.com/Mahaveer86619/Hearth/pkg/web"
)

func main() {
	config.LoadConfig()

	if err := db.InitRedis(config.AppConfig.RedisURL); err != nil {
		log.Fatalf("Critical: %v", err)
	}

	hub := services.NewHub()
	go hub.Run()

	pipe := ingestion_pipeline.NewIngestionPipelineService()
	pipe.Start()
	defer pipe.Stop()

	servers := web.NewServers()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, len(servers))

	// Start all servers
	for _, s := range servers {
		go func(srv web.Server) {
			logger.Info("Main", "Starting %s on %s", srv.Name(), srv.Addr())
			if err := srv.Start(); err != nil {
				errChan <- fmt.Errorf("%s error: %w", srv.Name(), err)
			}
		}(s)
	}

	select {
	case <-quit:
		logger.Info("Main", "Shutting down servers...")
	case err := <-errChan:
		logger.Error("Main", "Server startup failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown all servers
	for _, s := range servers {
		if err := s.Stop(ctx); err != nil {
			logger.Error("Main", "%s shutdown error: %v", s.Name(), err)
		}
	}

	logger.Info("Main", "Servers exited")
}
