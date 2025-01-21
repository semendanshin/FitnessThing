package repository

import (
	"context"
	"fitness-trainer/internal/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXRepository struct {
	pool *pgxpool.Pool
}

func NewPGXRepository(pool *pgxpool.Pool) *PGXRepository {
	return &PGXRepository{
		pool: pool,
	}
}

func (r *PGXRepository) runInTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}

	defer func () {
		if err := tx.Rollback(ctx); err != nil {
			logger.Errorf("failed to rollback transaction: %v", err)
		}
	}()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Errorf("failed to commit transaction: %v", err)
		return err
	}

	return nil
}

