package handler

import (
	"time"
	"net/http"

	"site-monitor/internal/domain"
)
// GetStatus godoc
// @Summary      Get site status
// @Description  Returns the latest monitoring status for a site
// @Tags         status
// @Produce      json
// @Param        id path string true "Site ID"
// @Success      200 {object} SiteStatusResponse
// @Failure      400 {string} string "Invalid site id"
// @Failure      404 {string} string "Site not found"
// @Router       /sites/{id}/status [get]
func (h *SiteHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
    siteID := r.PathValue("id")

	if siteID == "" {
		http.Error(w, "invalid site id", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	site, err := h.siteRepo.GetByID(ctx,siteID)
	if err != nil {
		http.Error(w, "site not found", http.StatusNotFound)
		return
	}

	status, ok := h.statusRepo.GetBySiteID(ctx,siteID)
	if !ok {
		resp := SiteStatusResponse{
			URL:    site.URL,
			Status: string(domain.StatusPending),
		}
		writeJSON(w, http.StatusOK, resp)
		return
	}

	var lastCheckedAt *time.Time

	if status.LastCheckedAt != nil {
    	lastCheckedAt = status.LastCheckedAt
	}

	var responseMs *int64
	if status.ResponseTime != nil && *status.ResponseTime > 0 {
		ms := status.ResponseTime.Milliseconds()
		responseMs = &ms
	}

	var statusCode *int
	if status.HTTPCode != 0 {
		code := status.HTTPCode
		statusCode = &code
	}

	var errMsg *string
	if status.Error != "" {
		err := status.Error
		errMsg = &err
	}

	resp := SiteStatusResponse{
		URL:           site.URL,
		Status:        string(status.Status),
		StatusCode:    statusCode,
		LastCheckedAt: lastCheckedAt,
		ResponseTime:  responseMs,
		Error:         errMsg,
	}

	writeJSON(w, http.StatusOK, resp)
}