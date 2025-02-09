package db

import (
	"context"
	"fitness-trainer/internal/logger"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContextKey string

const (
	EngineKey ContextKey = "db_engine"
)

type ContextManager struct {
	pool *pgxpool.Pool
}

func NewContextManager(pool *pgxpool.Pool) *ContextManager {
	return &ContextManager{
		pool: pool,
	}
}

type Engine interface {
	pgxscan.Querier
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type Transactioner interface {
	InTransaction(ctx context.Context, f func(ctx context.Context) error) error
}

func (cm *ContextManager) GetEngineFromContext(ctx context.Context) Engine {
	engine, ok := ctx.Value(EngineKey).(Engine)
	if !ok {
		return cm.pool
	}

	return engine
}

func (cm *ContextManager) PutEngineInContext(ctx context.Context, engine Engine) context.Context {
	return context.WithValue(ctx, EngineKey, engine)
}

func (cm *ContextManager) Begin(ctx context.Context) (context.Context, error) {
	_, ok := ctx.Value(EngineKey).(pgx.Tx)
	if ok {
		return ctx, nil
	}

	tx, err := cm.pool.Begin(ctx)
	if err != nil {
		return ctx, err
	}

	return cm.PutEngineInContext(ctx, tx), nil
}

func (cm *ContextManager) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(EngineKey).(*pgxpool.Tx)
	if !ok {
		return nil
	}

	return tx.Commit(ctx)
}

func (cm *ContextManager) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(EngineKey).(*pgxpool.Tx)
	if !ok {
		return nil
	}

	return tx.Rollback(ctx)
}

func (cm *ContextManager) InTransaction(ctx context.Context, f func(ctx context.Context) error) (err error) {
	txCtx, err := cm.Begin(ctx)
	if err != nil {
		return err
	}

	detCtx := context.WithoutCancel(txCtx)
	defer func() {
		if p := recover(); p != nil {
			logger.Errorf("panic occurred: %v", p)
			cm.Rollback(detCtx)
			panic(p)
		}
		if err != nil {
			logger.Errorf("error in tx occurred: %v", err)
			innerErr := cm.Rollback(txCtx)
			if innerErr != nil {
				logger.Errorf("failed to rollback transaction: %v", err)
			}
		} else {
			err = cm.Commit(txCtx)
			if err != nil {
				logger.Errorf("failed to commit transaction: %v", err)
			}
		}
	}()

	err = f(txCtx)

	return err
}

func (cm *ContextManager) Close() {
	cm.pool.Close()
}
