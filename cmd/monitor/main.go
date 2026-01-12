package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"site-monitor/internal/config"
	"site-monitor/internal/domain"
	"site-monitor/internal/http/handler"
	"site-monitor/internal/repository/memory"
	"site-monitor/internal/scheduler"
	"site-monitor/internal/server"
)
// @title           Site Monitor API
// @version         1.0
// @description     API for monitoring site availability
// @BasePath        /api/v1
func main() {
	configPath := flag.String("config", "", "Path to config YAML file")
	flag.Parse()

	if *configPath == "" {
		slog.Error("config path is required. Use -config <path>")
		os.Exit(1)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// =========================
	// Logger
	// =========================
	handlerJSON := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handlerJSON)

	logger.Info("Site Monitor started",
		slog.Int("num_sites", len(cfg.Sites)),
		slog.String("interval", cfg.Interval.String()),
	)

	// =========================
	// MAP config.Site -> domain.Site
	// =========================
	sites := make([]domain.Site, 0, len(cfg.Sites))
	for _, s := range cfg.Sites {
		sites = append(sites, domain.Site{
			ID:   s.ID,
			Name: s.Name,
			URL:  s.URL,
		})
	}

	// =========================
	// Repository + Handler
	// =========================
	// =========================

	siteRepo := memory.NewSiteMemoryRepository()
	siteStatus := memory.NewStatusMemoryRepository()

	siteRepo.Reset() 
	siteStatus.Reset()

	if err := memory.PopulateRepository(siteRepo, sites); err != nil {
    	logger.Error(fmt.Sprintf("failed to populate repository: %v", err))
	}
	startTime := time.Now()
	version := "1.0.0"

	siteHandler := handler.NewSiteHandler(siteRepo,siteStatus,logger)
	healthHandler := handler.NewHealthHandler(startTime, version)

	handlers := &server.Handlers{
		Site:   siteHandler,
		Health: healthHandler,
	}
	// =========================
	// Scheduler
	// =========================
	s := scheduler.New(cfg.Interval, cfg.Sites, logger, siteStatus)
	s.Start()

	// =========================
	// HTTP Router + Server
	// =========================
	router := server.NewRouter(handlers, logger)

	httpServer := server.New(":8080", router, logger)

	go func() {
		if err := httpServer.Start(); err != nil {
			logger.Error(
				"HTTP server error",
				slog.String("error", err.Error()),
			)
		}
	}()

	// =========================
	// Graceful shutdown
	// =========================
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Stop(ctx); err != nil {
		logger.Error("failed to stop HTTP server", slog.String("error", err.Error()))
	}

	s.Stop()

	logger.Info("Site Monitor stopped")
}
