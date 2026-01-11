package handler

import (
	"time"
	"net/http"

	"site-monitor/internal/domain"
)

func (h *SiteHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
    siteID := r.PathValue("id")

	if siteID == "" {
		http.Error(w, "invalid site id", http.StatusBadRequest)
		return
	}

	site, err := h.siteRepo.GetByID(siteID)
	if err != nil {
		http.Error(w, "site not found", http.StatusNotFound)
		return
	}

	status, ok := h.statusRepo.GetBySiteID(siteID)
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