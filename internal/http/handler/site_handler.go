package handler

import (
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
	repo repository.SiteRepository
}

func NewSiteHandler(repo repository.SiteRepository) *SiteHandler {
	return &SiteHandler{repo: repo}
}

type createSiteRequest struct {
	URL  string `json:"url"`
	Name string `json:"name,omitempty"`
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

	if err := h.repo.Create(site); err != nil {
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

	sites, err := h.repo.GetAll()
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
    vars := mux.Vars(r)      // получаем переменные из пути
    id := vars["id"]         // "id" — это имя параметра в роуте

    if err := h.repo.DeleteByID(id); err != nil {
        if err == repository.ErrSiteNotFound {
            http.Error(w, "site not found", http.StatusNotFound)
            return
        }
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
