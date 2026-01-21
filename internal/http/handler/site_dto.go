package handler

import "time"

type SiteStatusResponse struct {
	URL           string     `json:"url"`
	Status        string     `json:"status"`
	StatusCode    *int       `json:"statusCode,omitempty"`
	LastCheckedAt *time.Time `json:"lastCheckedAt,omitempty"`
	ResponseTime  *int64     `json:"responseTimeMs,omitempty"`
	Error         *string    `json:"error,omitempty"`
}

type createSiteRequest struct {
	URL  string `json:"url"`
	Name string `json:"name,omitempty"`
}
type SiteCheckHistoryItem struct {
    ID        string    `json:"id"`
    Status    string    `json:"status"`
    Code      int       `json:"code,omitempty"`
    Message   string    `json:"message,omitempty"`
    CheckedAt time.Time `json:"checked_at"`
}

// PaginatedResponse — универсальный ответ с пагинацией
type PaginatedResponse[T any] struct {
    Data  []T `json:"data"`
    Total int `json:"total"`
    Limit int `json:"limit"`
    Offset int `json:"offset"`
}