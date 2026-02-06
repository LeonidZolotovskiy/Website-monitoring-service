package repository

import (
	"log/slog"
	"context"
	"site-monitor/internal/domain"
)

type PostgresSiteAdapter struct {
	repo *PostgresSiteRepository
}

func NewPostgresSiteAdapter(repo *PostgresSiteRepository) *PostgresSiteAdapter {
	return &PostgresSiteAdapter{repo: repo}
}

func (a *PostgresSiteAdapter) Create(ctx context.Context, site domain.Site) error {
	_, err := a.repo.Create(ctx, site)
	return err
}

func (a *PostgresSiteAdapter) GetByID(ctx context.Context, id string) (*domain.Site, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *PostgresSiteAdapter) GetAll(ctx context.Context) ([]domain.Site, error) {
	return a.repo.GetAll(ctx)
}

func (a *PostgresSiteAdapter) Update(ctx context.Context, site domain.Site) error {
	return a.repo.Update(ctx, site)
}

func (a *PostgresSiteAdapter) DeleteByID(ctx context.Context, logger *slog.Logger,id string) error {
	return a.repo.DeleteByID(ctx, logger, id)
}

func (a *PostgresSiteAdapter) GetByURL(ctx context.Context, url string) (*domain.Site, error) {
	return a.repo.GetByURL(ctx, url)
}

func (a *PostgresSiteAdapter) GetHistoryBySiteID(ctx context.Context, siteID string, limit int, offset int) ([]domain.SiteCheckStatus, int, error) {
	return a.repo.GetHistoryBySiteID(ctx, siteID , limit, offset )
}