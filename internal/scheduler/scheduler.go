package scheduler

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"site-monitor/internal/checker"
	"site-monitor/internal/config"
)

type Scheduler struct {
	sites    []config.Site
	interval time.Duration
	logger   *slog.Logger
	client   *http.Client

	ticker  *time.Ticker
	quit    chan struct{}
	wg      sync.WaitGroup
}


func New(interval time.Duration, sites []config.Site, logger *slog.Logger) *Scheduler {
	if interval <= 0 {
		interval = time.Minute
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &Scheduler{
		sites:    sites,
		interval: interval,
		logger:   logger,
		client:   client,
		quit:     make(chan struct{}),
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
	var wg sync.WaitGroup

	for _, site := range s.sites {
		wg.Add(1)

		go func(site config.Site) {
			defer wg.Done()

			result := checker.CheckSite(s.client, site.URL)

			if result.Err != nil {
				s.logger.Warn("Site check failed",
					slog.String("url", site.URL),
					slog.String("name", site.Name),
					slog.String("error", result.Err.Error()),
				)
				return
			}

			if result.OK {
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
		}(site)
	}

	wg.Wait()
}


