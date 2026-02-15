package middleware

type ctxKey string

const RequestIDKey ctxKey = "request_id"

const RequestIDHeader = "X-Request-ID"
