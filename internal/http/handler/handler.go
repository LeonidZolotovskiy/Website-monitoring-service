package handler

import "site-monitor/internal/repository"

type SiteHandler struct {
	siteRepo   repository.SiteRepository
	statusRepo repository.StatusRepository
}

func NewSiteHandler(
	siteRepo repository.SiteRepository,
	statusRepo repository.StatusRepository,
) *SiteHandler {
	return &SiteHandler{
		siteRepo:   siteRepo,
		statusRepo: statusRepo,
	}
}
