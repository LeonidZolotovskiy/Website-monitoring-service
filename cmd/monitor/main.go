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

	// logger
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
	for i, s := range cfg.Sites {
		sites = append(sites, domain.Site{
			ID:   fmt.Sprintf("site-%d", i+1),
			Name: s.Name,
			URL:  s.URL,
		})
	}

	// =========================
	// Repository + Handler (DI)
	// =========================
	siteRepo := memory.NewSiteMemoryRepository(sites)
	siteHandler := handler.NewSiteHandler(siteRepo)

	// =========================
	// Scheduler (как было)
	// =========================
	s := scheduler.New(cfg.Interval, cfg.Sites, logger)
	s.Start()

	// =========================
	// HTTP server
	// =========================
	router := server.NewHTTPServer(siteHandler)
	httpServer := server.New(":8080", router, logger)
	httpServer.Start()

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
