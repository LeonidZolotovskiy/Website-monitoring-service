package postgres

import (
    "context"

    "log/slog"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type TxFunc func(tx pgx.Tx) error

func WithTransaction(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger, fn TxFunc) (err error) {
    tx, err := pool.Begin(ctx)
    if err != nil {
        return err
    }

    logger.Info("Transaction started")

    defer func() {
        if p := recover(); p != nil {
            _ = tx.Rollback(ctx)
            logger.Error("Transaction rolled back due to panic", slog.Any("panic", p))
            panic(p)
        } else if err != nil {
            _ = tx.Rollback(ctx)
            logger.Warn("Transaction rolled back due to error", slog.Any("err", err))
        } else {
            err = tx.Commit(ctx)
            if err != nil {
                logger.Error("Transaction commit failed", slog.Any("err", err))
            } else {
                logger.Info("Transaction committed successfully")
            }
        }
    }()

    err = fn(tx)
    return err
}
