package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(RequestIDHeader)
		if reqID == "" {
			reqID = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		r = r.WithContext(ctx)

		w.Header().Set(RequestIDHeader, reqID)

		next.ServeHTTP(w, r)
	})
}
