package server

import (
	"encoding/json"
	"net/http"
	"site-monitor/internal/http/handler"
)

type Handlers struct {
	Site   *handler.SiteHandler
	Health *handler.HealthHandler
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


