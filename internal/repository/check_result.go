package repository

import (
	"context"
	"site-monitor/internal/domain"
)

type CheckResultRepository interface {
	Create(ctx context.Context,result domain.CheckResult) error
	GetBySiteID(ctx context.Context,siteID string, limit, offset int) ([]domain.CheckResult, int, error)
	GetLatestBySiteID(ctx context.Context,siteID string) (*domain.CheckResult, error)
}
