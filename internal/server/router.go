package server

import (
	"net/http"
	"site-monitor/internal/http/handler"
)

func NewRouter(siteHandler *handler.SiteHandler) *http.ServeMux {
	root := http.NewServeMux()

	apiV1 := http.NewServeMux()

	apiV1.HandleFunc("/sites", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			siteHandler.GetSites(w, r)
		case http.MethodPost:
			siteHandler.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	apiV1.HandleFunc("/sites/", func(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		siteHandler.Delete(w, r)
		return
	}
	})
	apiV1.HandleFunc("/ping", PingHandler)

	root.Handle("/api/v1/", http.StripPrefix("/api/v1", apiV1))

	return root
}
