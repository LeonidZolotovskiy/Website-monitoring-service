package handler

import (
	"site-monitor/internal/repository"
	"log/slog"
)

type SiteHandler struct {
	siteRepo   repository.SiteRepository
	statusRepo repository.StatusRepository
	checkResultRepo repository.CheckResultRepository
	logger  *slog.Logger
}

func NewSiteHandler(
	siteRepo repository.SiteRepository,
	statusRepo repository.StatusRepository,
	checkResultRepo repository.CheckResultRepository,
	logger *slog.Logger,
) *SiteHandler {
	return &SiteHandler{
		siteRepo:   siteRepo,
		statusRepo: statusRepo,
		logger: logger,
	}
}
