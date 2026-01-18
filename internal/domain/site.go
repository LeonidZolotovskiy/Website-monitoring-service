package domain

import "time"

type Site struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	Name string `json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}