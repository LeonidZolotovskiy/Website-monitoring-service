package domain

import "time"

type SiteStatus string

const (
	StatusOK      SiteStatus = "ok"
	StatusError   SiteStatus = "error"
	StatusPending SiteStatus = "pending"
)

type SiteCheckStatus struct {
	SiteID        string        `json:"-"`
	StatusCode    *int
	URL           string        `json:"url"`
	Status        SiteStatus    `json:"status"`
	HTTPCode      int           `json:"http_code,omitempty"`
	LastCheckedAt *time.Time    `json:"last_checked_at,omitempty"`
	ResponseTime  *time.Duration `json:"response_time,omitempty"`
	Error         string        `json:"error,omitempty"`
}
