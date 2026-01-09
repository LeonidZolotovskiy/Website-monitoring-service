package server

import (
	"net/http"
	"site-monitor/internal/http/handler"
)

func NewRouter(siteHandler *handler.SiteHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/sites", siteHandler.GetSites)
	mux.HandleFunc("POST /api/v1/sites", siteHandler.Create)

	mux.HandleFunc("DELETE /api/v1/sites/{id}/", siteHandler.Delete)

	mux.HandleFunc("GET /api/v1/ping", PingHandler)

	mux.HandleFunc("GET /api/v1/sites/{id}/status", siteHandler.GetStatus)

	return mux
}
