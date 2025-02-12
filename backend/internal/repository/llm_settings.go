package repository

import (
	"context"
	"errors"
	"fitness-trainer/internal/domain"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/opentracing/opentracing-go"
)

type llmSettings struct {
	ID           pgtype.UUID
	UserID       pgtype.UUID
	BasePrompt   pgtype.Text
	VarietyLevel pgtype.Int8
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}

func (e llmSettings) toDomain() domain.GenerationSettings {
	return domain.GenerationSettings{
		Model: domain.Model{
			ID:        domain.ID(e.ID.Bytes),
			CreatedAt: e.CreatedAt.Time,
			UpdatedAt: e.UpdatedAt.Time,
		},
		UserID:       domain.ID(e.UserID.Bytes),
		BasePrompt:   e.BasePrompt.String,
		VarietyLevel: int(e.VarietyLevel.Int64),
	}
}

func llmSettingsFromDomain(settings domain.GenerationSettings) llmSettings {
	return llmSettings{
		ID:           uuidToPgtype(settings.ID),
		UserID:       uuidToPgtype(settings.UserID),
		BasePrompt:   pgtype.Text{String: settings.BasePrompt, Valid: true},
		VarietyLevel: pgtype.Int8{Int64: int64(settings.VarietyLevel), Valid: true},
		CreatedAt:    timeToPgtype(settings.CreatedAt),
		UpdatedAt:    timeToPgtype(settings.UpdatedAt),
	}
}

func (r *PGXRepository) GetGenerationSettings(ctx context.Context, userID domain.ID) (domain.GenerationSettings, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetGenerationSettings")
	defer span.Finish()

	const query = `
		SELECT id, user_id, base_prompt, variety_level, created_at, updated_at
		FROM llm_settings
		WHERE user_id = $1
	`

	engine := r.contextManager.GetEngineFromContext(ctx)

	var settings llmSettings
	if err := pgxscan.Get(ctx, engine, &settings, query, uuidToPgtype(userID)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.GenerationSettings{}, domain.ErrNotFound
		}
		return domain.GenerationSettings{}, err
	}

	return settings.toDomain(), nil
}

func (r *PGXRepository) SaveGenerationSettings(ctx context.Context, settings domain.GenerationSettings) (domain.GenerationSettings, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.SaveGenerationSettings")
	defer span.Finish()

	const query = `
		INSERT INTO llm_settings (id, user_id, base_prompt, variety_level, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE
		SET base_prompt = $3, variety_level = $4, updated_at = $6
		RETURNING id, user_id, base_prompt, variety_level, created_at, updated_at
	`

	engine := r.contextManager.GetEngineFromContext(ctx)

	settingsEntity := llmSettingsFromDomain(settings)

	if err := pgxscan.Get(ctx, engine, &settingsEntity, query,
		settingsEntity.ID, settingsEntity.UserID, settingsEntity.BasePrompt, settingsEntity.VarietyLevel, settingsEntity.CreatedAt, settingsEntity.UpdatedAt,
	); err != nil {
		return domain.GenerationSettings{}, err
	}

	return settingsEntity.toDomain(), nil
}
