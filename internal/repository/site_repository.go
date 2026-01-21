package repository

import (
	"errors"
	"site-monitor/internal/domain"
	"context"
)

var (
	ErrSiteAlreadyExists = errors.New("site already exists")
	ErrSiteNotFound      = errors.New("site not found")
	ErrDuplicateKey = errors.New("duplicate key")
)

type SiteRepository interface {
	Create(ctx context.Context, site domain.Site) error
	GetByURL(ctx context.Context, url string) (*domain.Site, error)
	GetAll(ctx context.Context) ([]domain.Site, error)
	DeleteByID(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*domain.Site, error)
}

type StatusRepository interface {
	Save(status domain.SiteCheckStatus)
	GetBySiteID(siteID string) (*domain.SiteCheckStatus, bool)
}