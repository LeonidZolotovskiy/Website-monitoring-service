package server

import (
	"net/http"
	"log/slog"

	"site-monitor/internal/http/middleware"

	_ "site-monitor/cmd/monitor/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(siteHandler *Handlers, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/health", siteHandler.Health.Check)
	mux.HandleFunc("GET /api/v1/ping", PingHandler)

	
	mux.HandleFunc("GET /api/v1/sites", siteHandler.Site.GetSites)
	mux.HandleFunc("POST /api/v1/sites", siteHandler.Site.Create)
	mux.HandleFunc("DELETE /api/v1/sites/{id}", siteHandler.Site.Delete)
	mux.HandleFunc("GET /api/v1/sites/{id}/status", siteHandler.Site.GetStatus)
	mux.HandleFunc("GET /api/v1/sites/{id}/history", siteHandler.Site.GetHistory)
	
	mux.Handle(
		"/swagger/swagger.json",
		http.StripPrefix(
			"/swagger/",
			http.FileServer(http.Dir("docs")),
		),
	)

	mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/swagger.json"),
		),
	)
	
	return middleware.Chain(
		mux,
		middleware.RequestID,
		middleware.Logging(logger),
		middleware.Recovery(logger),
	)
}
