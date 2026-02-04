package main

import (
	"context"
	"errors"
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
	"site-monitor/internal/repository"
	"site-monitor/internal/repository/memory"
	"site-monitor/internal/scheduler"
	"site-monitor/internal/server"
	"site-monitor/internal/storage/postgres"
)

func main() {
	configPath := flag.String("config", "", "Path to config YAML file")
	flag.Parse()
	handlerJSON := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handlerJSON)
	
	if *configPath == "" {
		slog.Error("config path is required. Use -config <path>")
		os.Exit(1)
	}
	
	cfg, err := config.Load(*configPath, logger)
	if err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// =========================
	// Logger
	// =========================
	
	

	logger.Info("Site Monitor started",
		slog.Int("num_sites", len(cfg.Sites)),
		slog.String("interval", cfg.Interval.String()),
	)

	// =========================
	// PostgreSQL pgx pool
	// =========================
	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg.DB)
	if err != nil {
		logger.Error("failed to connect to PostgreSQL", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("PostgreSQL connected (pgx pool)")
	if err := postgres.RunMigrations(pool, "migrations"); err != nil {
    	logger.Error("failed to run migrations", slog.String("error", err.Error()))
    	os.Exit(1)
	}
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
	// Repositories
	// =========================
	siteRepo := repository.NewPostgresSiteRepository(pool)
	//siteRepo := memory.NewSiteMemoryRepository()
	
	// Репозиторий истории проверок
	checkResultRepo := repository.NewPostgresCheckResultRepository(ctx, pool)

	siteStatus := memory.NewStatusMemoryRepository()
	siteStatus.Reset()

	// Populate initial sites из конфига
	for _, s := range sites {
		if _, err := siteRepo.Create(ctx, s); err != nil {
			if !errors.Is(err, repository.ErrDuplicateKey) {
				logger.Error(fmt.Sprintf("failed to populate site %s: %v", s.Name, err))
			}
		}
	}

	// =========================
	// Handlers
	// =========================
	startTime := time.Now()
	version := "1.0.0"

	siteHandler := handler.NewSiteHandler(siteRepo, siteStatus, checkResultRepo,logger)
	healthHandler := handler.NewHealthHandler(startTime, version)
	dbSites, err := siteRepo.GetAll(ctx)
	if err != nil {
		logger.Error("failed to load sites from DB", slog.String("error", err.Error()))
		os.Exit(1)
	}
	handlers := &server.Handlers{
		Site:   siteHandler,
		Health: healthHandler,
	}

	// =========================
	// Scheduler
	// =========================
	// Передаем репозиторий истории проверок в scheduler
	s := scheduler.New(cfg.Interval, dbSites, logger, siteStatus, checkResultRepo)
	s.Start()

	// =========================
	// HTTP Router + Server
	// =========================
	router := server.NewRouter(handlers, logger)
	httpServer := server.New(fmt.Sprintf(":%d", cfg.Port), router, logger)

	go func() {
		if err := httpServer.Start(); err != nil {
			logger.Error("HTTP server error", slog.String("error", err.Error()))
		}
	}()

	// =========================
	// Graceful shutdown
	// =========================
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("Shutting down...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Stop(ctxShutdown); err != nil {
		logger.Error("failed to stop HTTP server", slog.String("error", err.Error()))
	}

	s.Stop()

	if pool != nil {
		pool.Close()
	}

	logger.Info("Site Monitor stopped")
}
