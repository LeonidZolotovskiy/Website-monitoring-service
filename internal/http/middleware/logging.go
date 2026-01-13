package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			reqID, _ := r.Context().Value(RequestIDKey).(string)

			logger.Info(
				"incoming request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("request_id", reqID),
			)

			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			logger.Info(
				"request completed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rw.status),
				slog.Duration("duration", time.Since(start)),
				slog.String("request_id", reqID),
			)
		})
	}
}
