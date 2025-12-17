package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"site-monitor/internal/config"
	"site-monitor/internal/scheduler"
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

	// Создаём JSON-логгер
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)

	logger.Info("Site Monitor started",
		slog.Int("num_sites", len(cfg.Sites)),
		slog.String("interval", cfg.Interval.String()),
	)

	s := scheduler.New(cfg.Interval, cfg.Sites, logger)
	s.Start()

	// Обработка Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("Shutting down...")
	s.Stop()
	logger.Info("Site Monitor stopped")
}
