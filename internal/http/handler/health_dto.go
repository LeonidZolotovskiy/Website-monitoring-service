package handler

import "time"

type HealthStatus string

const (
	HealthHealthy   HealthStatus = "healthy"
	HealthUnhealthy HealthStatus = "unhealthy"
)

type DependencyStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type HealthResponse struct {
	Status       HealthStatus        `json:"status"`
	Version      string              `json:"version"`
	Uptime       string              `json:"uptime"`
	ServerTime   time.Time           `json:"serverTime"`
	Dependencies []DependencyStatus  `json:"dependencies,omitempty"`
}
