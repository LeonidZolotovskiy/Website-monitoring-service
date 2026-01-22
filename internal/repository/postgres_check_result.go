package repository

import (
	"context"
	"fmt"
	"site-monitor/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCheckResultRepository struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewPostgresCheckResultRepository(ctx context.Context, pool *pgxpool.Pool) *PostgresCheckResultRepository {
	return &PostgresCheckResultRepository{pool: pool, ctx: ctx}
}

func (r *PostgresCheckResultRepository) Create(result domain.CheckResult) error {
	query := `INSERT INTO site_checks (site_id, status, response_time_ms, checked_at) VALUES ($1, $2, $3, NOW());`
	_, err := r.pool.Exec(r.ctx, query, result.SiteID, result.Status, result.ResponseTimeMs)
	if err != nil {
		return fmt.Errorf("insert check result: %w", err)
	}
	return err
}

func (r *PostgresCheckResultRepository) GetBySiteID(siteID string, limit, offset int) ([]domain.CheckResult, error) {
	query := `SELECT id, site_id, status, response_time_ms, checked_at FROM site_checks WHERE site_id = $1 ORDER BY checked_at DESC LIMIT $2 OFFSET $3;`
	rows, err := r.pool.Query(r.ctx, query, siteID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.CheckResult
	for rows.Next() {
		var rsl domain.CheckResult
		if err := rows.Scan(&rsl.ID, &rsl.SiteID, &rsl.Status, &rsl.ResponseTimeMs, &rsl.CheckedAt); err != nil {
			continue
		}
		results = append(results, rsl)
	}
	return results, nil
}

func (r *PostgresCheckResultRepository) GetLatestBySiteID(siteID string) (*domain.CheckResult, error) {
	query := `SELECT id, site_id, status, response_time_ms, checked_at FROM site_checks WHERE site_id = $1 ORDER BY checked_at DESC LIMIT 1;`
	var rsl domain.CheckResult
	err := r.pool.QueryRow(r.ctx, query, siteID).Scan(&rsl.ID, &rsl.SiteID, &rsl.Status, &rsl.ResponseTimeMs, &rsl.CheckedAt)
	if err != nil {
		return nil, err
	}
	return &rsl, nil
}
