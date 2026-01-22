package scheduler

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"site-monitor/internal/domain"
	"site-monitor/internal/checker"
	"site-monitor/internal/repository"
)

type StatusRepository interface {
	Save(status domain.SiteCheckStatus)
}

type Scheduler struct {
	sites    []domain.Site
	interval time.Duration
	logger   *slog.Logger
	client   *http.Client
	statusRepo StatusRepository

	checkResultRepo repository.CheckResultRepository

	ticker  *time.Ticker
	quit    chan struct{}
	wg      sync.WaitGroup
}


func New(
	interval time.Duration,
	sites []domain.Site,
	logger *slog.Logger,
	statusRepo StatusRepository,
	checkResultRepo repository.CheckResultRepository,
) *Scheduler {
	if interval <= 0 {
		interval = time.Minute
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &Scheduler{
		sites:      sites,
		interval:   interval,
		logger:     logger,
		client:     client,
		statusRepo: statusRepo,
		quit:       make(chan struct{}),
		checkResultRepo: checkResultRepo,
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

		go func(site domain.Site) {
			defer wg.Done()

			start := time.Now()
			result := checker.CheckSite(s.client, site.URL)
			checkedAt := time.Now()
			responseTime := checkedAt.Sub(start)

			status := domain.StatusOK
			if !result.OK || result.Err != nil {
				status = domain.StatusError
			}

			var statusCode *int
			if result.Err == nil {
				statusCode = &result.StatusCode
			}

			s.statusRepo.Save(domain.SiteCheckStatus{
				SiteID:       site.ID, 
				Status:       status,
				StatusCode:   statusCode,
				LastCheckedAt:    &checkedAt,
				ResponseTime: &responseTime,
			})

			// логирование — как было
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



