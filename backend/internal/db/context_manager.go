package db

import (
	"context"

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

func (cm *ContextManager) Close() {
	cm.pool.Close()
}
