package postgres

import (
    "database/sql"
    "fmt"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"

    _ "github.com/golang-migrate/migrate/v4/source/file" 
    _ "github.com/jackc/pgx/v5/stdlib"                  
)

func RunMigrations(connString, path string) error {
    db, err := sql.Open("pgx", connString) // use pgx driver
    if err != nil {
        return fmt.Errorf("open sql db: %w", err)
    }
    defer db.Close()

    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return fmt.Errorf("create migrate driver: %w", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://"+path,
        "postgres",
        driver,
    )
    if err != nil {
        return fmt.Errorf("create migrate instance: %w", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("run migrations: %w", err)
    }

    return nil
}

