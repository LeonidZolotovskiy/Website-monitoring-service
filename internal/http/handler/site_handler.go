package handler

import (
	"encoding/json"
	"net/http"

	"site-monitor/internal/repository"
)

type SiteHandler struct {
	repo repository.SiteRepository
}

func NewSiteHandler(repo repository.SiteRepository) *SiteHandler {
	return &SiteHandler{repo: repo}
}

func (h *SiteHandler) GetSites(w http.ResponseWriter, r *http.Request) {
	sites := h.repo.GetAll()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(sites)
}
