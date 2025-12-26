package repository

import "site-monitor/internal/domain"

type SiteRepository interface {
	GetAll() []domain.Site
}