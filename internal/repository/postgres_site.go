package repository

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"site-monitor/internal/domain"
	"site-monitor/internal/storage/postgres"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type PostgresSiteRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresSiteRepository(pool *pgxpool.Pool) *PostgresSiteRepository {
	return &PostgresSiteRepository{pool: pool}
}


func (r *PostgresSiteRepository) Create(ctx context.Context, site domain.Site) (string, error) {
	query := `
		INSERT INTO sites (name, url, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id;
	`
	var id string
	err := r.pool.QueryRow(ctx, query, site.Name, site.URL).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // уникальный ключ
			return "", ErrDuplicateKey
		}
		log.Printf("failed to create site: %v", err)
		return "", err
	}
	return id, nil
}

func (r *PostgresSiteRepository) GetByID(ctx context.Context, id string) (*domain.Site, error) {
	query := `
		SELECT id, name, url, created_at, updated_at
		FROM sites
		WHERE id = $1;
	`
	site := &domain.Site{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&site.ID, &site.Name, &site.URL, &site.CreatedAt, &site.UpdatedAt,
	)
	if err != nil {
		return nil, ErrSiteNotFound
	}
	return site, nil
}

func (r *PostgresSiteRepository) GetAll(ctx context.Context) ([]domain.Site, error) {
	query := `
		SELECT id, name, url, created_at, updated_at
		FROM sites;
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		log.Printf("failed to get all sites: %v", err)
		return nil, err
	}
	defer rows.Close()

	var sites []domain.Site
	for rows.Next() {
		var s domain.Site
		err := rows.Scan(&s.ID, &s.Name, &s.URL, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			log.Printf("failed to scan site: %v", err)
			continue
		}
		sites = append(sites, s)
	}
	return sites, nil
}

func (r *PostgresSiteRepository) Update(ctx context.Context, site domain.Site) error {
	query := `
		UPDATE sites
		SET name = $1, url = $2, updated_at = NOW()
		WHERE id = $3;
	`
	cmdTag, err := r.pool.Exec(ctx, query, site.Name, site.URL, site.ID)
	if err != nil {
		log.Printf("failed to update site: %v", err)
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrSiteNotFound
	}
	return nil
}

func (r *PostgresSiteRepository) Delete(ctx context.Context,logger *slog.Logger ,siteID string) error {
    return postgres.WithTransaction(ctx, r.pool, logger, func(tx pgx.Tx) error {

        if _, err := tx.Exec(ctx, `DELETE FROM site_checks WHERE site_id = $1`, siteID); err != nil {
            return err
        }

        cmdTag, err := tx.Exec(ctx, `DELETE FROM sites WHERE id = $1`, siteID)
        if err != nil {
            return err
        }
        if cmdTag.RowsAffected() == 0 {
            return ErrSiteNotFound
        }
        return nil
    })
}

func (r *PostgresSiteRepository) GetByURL(ctx context.Context, url string) (*domain.Site, error) {
	query := `
		SELECT id, name, url, created_at, updated_at
		FROM sites
		WHERE url = $1;
	`
	site := &domain.Site{}
	err := r.pool.QueryRow(ctx, query, url).Scan(
		&site.ID, &site.Name, &site.URL, &site.CreatedAt, &site.UpdatedAt,
	)
	if err != nil {
		return nil, ErrSiteNotFound
	}
	return site, nil
}

func (r *PostgresSiteRepository) GetHistoryBySiteID(ctx context.Context, siteID string, limit, offset int) ([]domain.SiteCheckStatus, int, error) {
    var total int
    err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM site_checks WHERE site_id = $1`, siteID).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    query := `
        SELECT id, status, code, message, checked_at
        FROM site_checks
        WHERE site_id = $1
        ORDER BY checked_at DESC
        LIMIT $2 OFFSET $3
    `
    rows, err := r.pool.Query(ctx, query, siteID, limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var history []domain.SiteCheckStatus
    for rows.Next() {
        var s domain.SiteCheckStatus
        var status string
        var code *int64   
        var msg *string

        if err := rows.Scan(&s.SiteID, &status, &code, &msg, &s.LastCheckedAt); err != nil {
            continue
        }

        s.Status = domain.SiteStatus(status)

        if code != nil {
            v := int(*code)   
            s.StatusCode = &v
        }

        history = append(history, s)
    }

    return history, total, nil
}
