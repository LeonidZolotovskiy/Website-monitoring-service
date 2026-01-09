package handler

import (
	"time"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"site-monitor/internal/domain"
	"site-monitor/internal/repository"
)

type SiteHandler struct {
	siteRepo   repository.SiteRepository
	statusRepo repository.StatusRepository
}

type SiteStatusResponse struct {
	URL           string  `json:"url"`
	Status        string  `json:"status"`
	StatusCode    *int    `json:"statusCode,omitempty"`
	LastCheckedAt *string `json:"lastCheckedAt,omitempty"`
	ResponseTime  *int64  `json:"responseTimeMs,omitempty"`
	Error         *string `json:"error,omitempty"`
}

type createSiteRequest struct {
	URL  string `json:"url"`
	Name string `json:"name,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

	if err := h.siteRepo.Create(site); err != nil {
		if err == repository.ErrSiteAlreadyExists {
			http.Error(w, "site with this url already exists", http.StatusConflict)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(site)
}

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

func (h *SiteHandler) Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)      
    id := vars["id"]        

    if err := h.siteRepo.DeleteByID(id); err != nil {
        if err == repository.ErrSiteNotFound {
            http.Error(w, "site not found", http.StatusNotFound)
            return
        }
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *SiteHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)      
    siteID := vars["id"]    

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

	var lastCheckedAt *string
	if status.LastCheckedAt != nil {
		t := status.LastCheckedAt.Format(time.RFC3339)
		lastCheckedAt = &t
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

