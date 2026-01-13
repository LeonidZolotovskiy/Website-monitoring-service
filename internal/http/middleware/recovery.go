package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					reqID, _ := r.Context().Value(RequestIDKey).(string)

					logger.Error(
						"panic recovered",
						slog.Any("error", err),
						slog.String("request_id", reqID),
						slog.String("stacktrace", string(debug.Stack())),
					)

					http.Error(
						w,
						http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError,
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
