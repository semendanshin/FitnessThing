package repository

import (
	"fitness-trainer/internal/db"
)

type PGXRepository struct {
	contextManager *db.ContextManager
}

func NewPGXRepository(ctxManager *db.ContextManager) *PGXRepository {
	return &PGXRepository{
		contextManager: ctxManager,
	}
}

// func (r *PGXRepository) runInTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
// 	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
// 		IsoLevel: pgx.ReadCommitted,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	defer func() {
// 		if err := tx.Rollback(ctx); err != nil {
// 			logger.Errorf("failed to rollback transaction: %v", err)
// 		}
// 	}()

// 	if err := fn(tx); err != nil {
// 		return err
// 	}

// 	if err := tx.Commit(ctx); err != nil {
// 		logger.Errorf("failed to commit transaction: %v", err)
// 		return err
// 	}

// 	return nil
// }
