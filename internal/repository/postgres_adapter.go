package repository

import (
	"context"
	"site-monitor/internal/domain"
)

type PostgresSiteAdapter struct {
	repo *PostgresSiteRepository
	ctx  context.Context
}

func NewPostgresSiteAdapter(ctx context.Context, repo *PostgresSiteRepository) *PostgresSiteAdapter {
	return &PostgresSiteAdapter{repo: repo, ctx: ctx}
}

func (a *PostgresSiteAdapter) Create(site domain.Site) error {
	_, err := a.repo.Create(a.ctx, site)
	return err
}

func (a *PostgresSiteAdapter) GetByID(id string) (*domain.Site, error) {
	return a.repo.GetByID(a.ctx, id)
}

func (a *PostgresSiteAdapter) GetAll() ([]domain.Site, error) {
	return a.repo.GetAll(a.ctx)
}

func (a *PostgresSiteAdapter) Update(site domain.Site) error {
	return a.repo.Update(a.ctx, site)
}

func (a *PostgresSiteAdapter) DeleteByID(id string) error {
	return a.repo.Delete(a.ctx, id)
}

func (a *PostgresSiteAdapter) GetByURL(url string) (*domain.Site, error) {
	return a.repo.GetByURL(a.ctx, url)
}