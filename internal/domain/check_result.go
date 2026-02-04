package domain

import "time"

type CheckResult struct {
	ID           string    `json:"id"`
    IsAvailable  bool      `json:"is_available"`
    HTTPStatus   *int      `json:"http_status,omitempty"`
    ResponseTime *time.Duration `json:"response_time_ms,omitempty"`
    ErrorMessage *string   `json:"error_message,omitempty"`
    CheckedAt    time.Time `json:"checked_at"`
}

type PaginatedResponse[T any] struct {
    Data  []T `json:"data"`
    Total int `json:"total"`
    Limit int `json:"limit"`
    Offset int `json:"offset"`
}