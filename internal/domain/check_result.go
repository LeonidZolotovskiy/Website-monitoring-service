package domain

import "time"

type CheckResult struct {
	ID             string    `json:"id"`              // UUID или SERIAL
	SiteID         string    `json:"site_id"`         // ID сайта
	Status         string    `json:"status"`          // "up" или "down"
	ResponseTimeMs int       `json:"response_time_ms"`// Время отклика в миллисекундах
	CheckedAt      time.Time `json:"checked_at"`      // Время проверки
}
