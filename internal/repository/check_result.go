package repository

import (
	"site-monitor/internal/domain"
)

type CheckResultRepository interface {
	Create(result domain.CheckResult) error
	GetBySiteID(siteID string, limit, offset int) ([]domain.CheckResult, error)
	GetLatestBySiteID(siteID string) (*domain.CheckResult, error)
}
