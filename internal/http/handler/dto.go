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
