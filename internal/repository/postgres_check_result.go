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

func (r *PostgresCheckResultRepository) Create(ctx context.Context, result domain.CheckResult) error {
	query := `
		INSERT INTO site_checks (site_id, http_status, response_time_ms, checked_at, is_available)
		VALUES ($1, $2, $3, NOW());
	`

	_, err := r.pool.Exec(ctx, query, result.HTTPStatus, result.ResponseTime, result.CheckedAt, result.IsAvailable)
	if err != nil {
		return fmt.Errorf("insert check result: %w", err)
	}

	return nil
}


func (r *PostgresCheckResultRepository) GetBySiteID(
	ctx context.Context,
	siteID string,
	limit, offset int,
) ([]domain.CheckResult, int, error) {

	var total int

	countQuery := `SELECT COUNT(*) FROM site_checks WHERE site_id = $1;`
	err := r.pool.QueryRow(ctx, countQuery, siteID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, site_id, http_status, response_time_ms, checked_at
		FROM site_checks
		WHERE site_id = $1
		ORDER BY checked_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.pool.Query(ctx, query, siteID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []domain.CheckResult

	for rows.Next() {
		var rsl domain.CheckResult

		if err := rows.Scan(
			&rsl.ID,
			&rsl.HTTPStatus,
			&rsl.IsAvailable,
			&rsl.ResponseTime,
			&rsl.CheckedAt,
		); err != nil {
			return nil, 0, err
		}

		results = append(results, rsl)
	}

	return results, total, rows.Err()
}


func (r *PostgresCheckResultRepository) GetLatestBySiteID(
	ctx context.Context,
	siteID string,
) (*domain.CheckResult, error) {

	query := `
		SELECT id, site_id, http_status, response_time_ms, checked_at
		FROM site_checks
		WHERE site_id = $1
		ORDER BY checked_at DESC
		LIMIT 1;
	`

	var rsl domain.CheckResult

	err := r.pool.QueryRow(ctx, query, siteID).Scan(
		&rsl.ID,
		&rsl.HTTPStatus,
		&rsl.IsAvailable,
		&rsl.ResponseTime,
		&rsl.CheckedAt,
	)

	if err != nil {
		return nil, err
	}

	return &rsl, nil
}

