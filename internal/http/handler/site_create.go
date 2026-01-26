package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"

	"site-monitor/internal/domain"
	"site-monitor/internal/repository"
)
// Create godoc
// @Summary      Create site
// @Description  Creates a new site for monitoring
// @Tags         sites
// @Accept       json
// @Produce      json
// @Param        request body createSiteRequest true "Site data"
// @Success      201 {object} domain.Site "Site created"
// @Failure      400 {string} string "Invalid request body or URL"
// @Failure      409 {string} string "Site already exists"
// @Failure      500 {string} string "Internal server error"
// @Router       /sites [post]
func (h *SiteHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req createSiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.URL) == "" {
		http.Error(w, "url must not be empty", http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		http.Error(w, "invalid url format", http.StatusBadRequest)
		return
	}

	site := domain.Site{
		ID:   uuid.NewString(),
		URL:  req.URL,
		Name: req.Name,
	}
	ctx := r.Context()
	if  _,err := h.siteRepo.Create(ctx,site); err != nil {
		if err == repository.ErrSiteAlreadyExists {
			http.Error(w, "site with this url already exists", http.StatusConflict)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, site)
}