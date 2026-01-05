package server

import (
	"encoding/json"
	"net/http"
	"site-monitor/internal/http/handler"
	"site-monitor/internal/repository"
)

type Handlers struct {
	Site *handler.SiteHandler
}

func NewHandlers(siteRepo repository.SiteRepository) *Handlers {
	return &Handlers{
		Site: handler.NewSiteHandler(siteRepo),
	}
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "pong",
	})
}

