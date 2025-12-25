package server

import "net/http"

func NewRouter() *http.ServeMux {
	root := http.NewServeMux()

	apiV1 := http.NewServeMux()
	apiV1.HandleFunc("/ping", PingHandler)

	root.Handle("/api/v1/", http.StripPrefix("/api/v1", apiV1))

	return root
}
