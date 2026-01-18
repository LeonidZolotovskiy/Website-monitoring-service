package repository

import (
	"errors"
	"site-monitor/internal/domain"
)

var (
	ErrSiteAlreadyExists = errors.New("site already exists")
	ErrSiteNotFound      = errors.New("site not found")
	ErrDuplicateKey = errors.New("duplicate key")
)

type SiteRepository interface {
	Create(site domain.Site) error
	GetByURL(url string) (*domain.Site, error)
	GetAll() ([]domain.Site, error)
	DeleteByID(id string) error
	GetByID(id string) (*domain.Site, error)
}

type StatusRepository interface {
	Save(status domain.SiteCheckStatus)
	GetBySiteID(siteID string) (*domain.SiteCheckStatus, bool)
}