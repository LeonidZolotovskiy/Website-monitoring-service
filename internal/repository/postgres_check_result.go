package repository

import (
	"context"
	"fmt"
	"site-monitor/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCheckResultRepository struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewPostgresCheckResultRepository(ctx context.Context, pool *pgxpool.Pool) *PostgresCheckResultRepository {
	return &PostgresCheckResultRepository{pool: pool, ctx: ctx}
}

func (r *PostgresCheckResultRepository) Create(
	ctx context.Context,
	result domain.CheckResult,
) error {

	if result.SiteID == "" {
		return fmt.Errorf("cannot insert check result: SiteID is empty")
	}

	const query = `
        INSERT INTO site_checks (
            site_id,
            http_status,
            response_time_ms,
            is_available,
            error_message,
            checked_at
        )
        VALUES ($1,$2,$3,$4,$5,$6);
    `

	var responseTime *int32
	if result.ResponseTime != nil {
		rt := int32(*result.ResponseTime / time.Millisecond)
		responseTime = &rt
	}

	_, err := r.pool.Exec(
		ctx,
		query,
		result.SiteID,
		result.HTTPStatus,
		responseTime,
		result.IsAvailable,
		result.ErrorMessage,
		result.CheckedAt,
	)

	if err != nil {
		return fmt.Errorf("insert check result: %w", err)
	}

	return nil
}



// GetBySiteID возвращает список проверок сайта с пагинацией
func (r *PostgresCheckResultRepository) GetBySiteID(
	ctx context.Context,
	siteID string,
	limit, offset int,
) ([]domain.CheckResult, int, error) {

	var total int
	countQuery := `SELECT COUNT(*) FROM site_checks WHERE site_id = $1;`
	if err := r.pool.QueryRow(ctx, countQuery, siteID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count check results: %w", err)
	}

	query := `
		SELECT id, site_id, http_status, response_time_ms, is_available, error_message, checked_at
		FROM site_checks
		WHERE site_id = $1
		ORDER BY checked_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.pool.Query(ctx, query, siteID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query check results: %w", err)
	}
	defer rows.Close()

	var results []domain.CheckResult

	for rows.Next() {
		var rsl domain.CheckResult
		var responseTimeMs *int32

		if err := rows.Scan(
			&rsl.ID,
			&rsl.SiteID,
			&rsl.HTTPStatus,
			&responseTimeMs,
			&rsl.IsAvailable,
			&rsl.ErrorMessage,
			&rsl.CheckedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan row: %w", err)
		}

		if responseTimeMs != nil {
			rt := time.Duration(*responseTimeMs) * time.Millisecond
			rsl.ResponseTime = &rt
		}

		results = append(results, rsl)
	}

	return results, total, rows.Err()
}

// GetLatestBySiteID возвращает последнюю проверку сайта
func (r *PostgresCheckResultRepository) GetLatestBySiteID(
	ctx context.Context,
	siteID string,
) (*domain.CheckResult, error) {

	query := `
		SELECT id, site_id, http_status, response_time_ms, is_available, error_message, checked_at
		FROM site_checks
		WHERE site_id = $1
		ORDER BY checked_at DESC
		LIMIT 1;
	`

	var rsl domain.CheckResult
	var responseTimeMs *int32

	err := r.pool.QueryRow(ctx, query, siteID).Scan(
		&rsl.ID,
		&rsl.SiteID,
		&rsl.HTTPStatus,
		&responseTimeMs,
		&rsl.IsAvailable,
		&rsl.ErrorMessage,
		&rsl.CheckedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get latest check result: %w", err)
	}

	if responseTimeMs != nil {
		rt := time.Duration(*responseTimeMs) * time.Millisecond
		rsl.ResponseTime = &rt
	}

	return &rsl, nil
}
