package handler

import (
	"encoding/json"
	"net/http"
)
// GetSites godoc
// @Summary      Get sites
// @Description  Returns a list of all monitored sites
// @Tags         sites
// @Produce      json
// @Success      200 {array} domain.Site
// @Failure      500 {string} string "Failed to get sites"
// @Router       /sites [get]
func (h *SiteHandler) GetSites(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sites, err := h.siteRepo.GetAll()
	if err != nil {
		http.Error(w, "failed to get sites", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(sites); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}