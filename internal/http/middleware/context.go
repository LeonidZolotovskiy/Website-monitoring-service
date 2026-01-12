package middleware

type ctxKey string

const (
	RequestIDKey ctxKey = "request_id"
	RequestIDHeader     = "X-Request-ID"
)
