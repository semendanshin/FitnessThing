package repository

import (
	"context"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/logger"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/opentracing/opentracing-go"
)

type sessionEntity struct {
	ID        pgtype.UUID
	UserID    pgtype.UUID
	ExpiredAt pgtype.Timestamptz
	Token     string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

func (s sessionEntity) toDomain() domain.Session {
	return domain.Session{
		Model: domain.Model{
			ID:        domain.ID(s.ID.Bytes),
			CreatedAt: s.CreatedAt.Time,
			UpdatedAt: s.UpdatedAt.Time,
		},
		UserID:    domain.ID(s.UserID.Bytes),
		ExpiredAt: s.ExpiredAt.Time,
		Token:     s.Token,
	}
}

func (r *PGXRepository) GetSessionByToken(ctx context.Context, token string) (domain.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetSessionByToken")
	defer span.Finish()

	query := `
		SELECT * FROM sessions s
		WHERE s.token = $1
	`

	var session sessionEntity
	err := pgxscan.Get(ctx, r.pool, &session, query, token)
	if err != nil {
		logger.Errorf("failed to get session by token: %v", err)
		if err == pgx.ErrNoRows {
			return domain.Session{}, domain.ErrNotFound
		}
		return domain.Session{}, err
	}

	return session.toDomain(), nil
}

func (r *PGXRepository) SetSessionExpired(ctx context.Context, id domain.ID, expiredAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.SetSessionExpired")
	defer span.Finish()

	query := `
		UPDATE sessions s
		SET expired_at = $2
		WHERE id = $1
		RETURNING id
	`

	var session sessionEntity
	err := pgxscan.Get(ctx, r.pool, &session, query, uuidToPgtype(id), timeToPgtype(expiredAt))
	if err != nil {
		logger.Errorf("failed to set session expired: %v", err)
		return err
	}

	return nil
}

func (r *PGXRepository) CreateSession(ctx context.Context, session domain.Session) (domain.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.CreateSession")
	defer span.Finish()

	query := `
		INSERT INTO sessions (id, user_id, token, expired_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`

	var sessionEntity sessionEntity
	err := pgxscan.Get(
		ctx,
		r.pool,
		&sessionEntity,
		query,
		uuidToPgtype(session.ID),
		uuidToPgtype(session.UserID),
		session.Token,
		timeToPgtype(session.ExpiredAt),
		timeToPgtype(session.CreatedAt),
	)
	if err != nil {
		logger.Errorf("failed to create session: %v", err)
		return domain.Session{}, err
	}

	return sessionEntity.toDomain(), nil
}
