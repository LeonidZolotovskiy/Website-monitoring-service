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

func (s *Server) Start() {
	go func() {
		s.logger.Info("HTTP server started", slog.String("addr", s.httpServer.Addr))

		if err := s.httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("HTTP server error", slog.String("error", err.Error()))
		}
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
