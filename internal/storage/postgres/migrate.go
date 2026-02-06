package postgres

import (
    "errors"
    "fmt"


    
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    "github.com/jackc/pgx/v5/pgxpool"

    _ "github.com/golang-migrate/migrate/v4/source/file" 
    "github.com/jackc/pgx/v5/stdlib"                  
)

func RunMigrations(pool *pgxpool.Pool, path string) error {
    config := pool.Config().ConnConfig

    db := stdlib.OpenDB(*config)

    // важно — ограничиваем коннекты
    db.SetMaxOpenConns(1)

    defer db.Close()

    if err := db.Ping(); err != nil {
        return fmt.Errorf("ping db: %w", err)
    }

    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return fmt.Errorf("create migrate driver: %w", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        fmt.Sprintf("file://%s", path),
        "postgres",
        driver,
    )
    
    if err != nil {
        return fmt.Errorf("create migrate instance: %w", err)
    }

    if err := m.Up(); err != nil {
        if errors.Is(err, migrate.ErrNoChange) {
            return nil
        }

        var dirty migrate.ErrDirty
        if errors.As(err, &dirty) {
            return fmt.Errorf("database is dirty at version %d", dirty.Version)
        }

        return fmt.Errorf("run migrations: %w", err)
    }

    return nil
}
