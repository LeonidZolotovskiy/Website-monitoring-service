package scheduler

import (
	"log/slog"
	"sync"
	"time"

	"site-monitor/internal/checker"
	"site-monitor/internal/config"
)

type Scheduler struct {
	sites    []config.Site
	interval time.Duration
	logger   *slog.Logger

	ticker  *time.Ticker
	quit    chan struct{}
	wg      sync.WaitGroup
	mu      sync.Mutex
	running bool
}

func New(interval time.Duration, sites []config.Site, logger *slog.Logger) *Scheduler {
	if interval <= 0 {
		interval = time.Minute
	}
	return &Scheduler{
		sites:    sites,
		interval: interval,
		quit:     make(chan struct{}),
		logger:   logger,
	}
}

func (s *Scheduler) Start() {
	s.ticker = time.NewTicker(s.interval)
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		s.runChecks() // первая проверка сразу

		for {
			select {
			case <-s.ticker.C:
				s.runChecks()
			case <-s.quit:
				s.logger.Info("Scheduler stopped")
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.quit)
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.wg.Wait()
}

func (s *Scheduler) runChecks() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	for _, site := range s.sites {
		result := checker.CheckSite(site.URL)

		if result.Err != nil {
			s.logger.Warn("Site check failed",
				slog.String("url", site.URL),
				slog.String("name", site.Name),
				slog.String("error", result.Err.Error()),
			)
		} else if result.OK {
			s.logger.Info("Site is up",
				slog.String("url", site.URL),
				slog.String("name", site.Name),
				slog.Int("status_code", result.StatusCode),
			)
		} else {
			s.logger.Warn("Site is down",
				slog.String("url", site.URL),
				slog.String("name", site.Name),
				slog.Int("status_code", result.StatusCode),
			)
		}
	}

	s.mu.Lock()
	s.running = false
	s.mu.Unlock()
}

