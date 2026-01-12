package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func New(addr string, handler http.Handler, logger *slog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info("HTTP server started", slog.String("addr", s.httpServer.Addr))

	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server...")
	return s.httpServer.Shutdown(ctx)
}