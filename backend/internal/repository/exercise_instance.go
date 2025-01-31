package repository

import (
	"context"

	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/logger"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/opentracing/opentracing-go"
)

type exerciseInstanceEntity struct {
	ID         pgtype.UUID
	RoutineID  pgtype.UUID
	ExerciseID pgtype.UUID
	CreatedAt  pgtype.Timestamptz
	UpdatedAt  pgtype.Timestamptz
}

func (e exerciseInstanceEntity) toDomain() domain.ExerciseInstance {
	return domain.ExerciseInstance{
		Model: domain.Model{
			ID:        domain.ID(e.ID.Bytes),
			CreatedAt: e.CreatedAt.Time,
			UpdatedAt: e.UpdatedAt.Time,
		},
		RoutineID:  domain.ID(e.RoutineID.Bytes),
		ExerciseID: domain.ID(e.ExerciseID.Bytes),
	}
}

func exerciseInstanceFromDomain(exerciseInstance domain.ExerciseInstance) exerciseInstanceEntity {
	return exerciseInstanceEntity{
		ID:         uuidToPgtype(exerciseInstance.ID),
		RoutineID:  uuidToPgtype(exerciseInstance.RoutineID),
		ExerciseID: uuidToPgtype(exerciseInstance.ExerciseID),
		CreatedAt:  timeToPgtype(exerciseInstance.CreatedAt),
		UpdatedAt:  timeToPgtype(exerciseInstance.UpdatedAt),
	}
}

func (r *PGXRepository) GetExerciseInstanceByID(ctx context.Context, id domain.ID) (domain.ExerciseInstance, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetExerciseInstanceByID")
	defer span.Finish()

	query := `
		SELECT * FROM exercise_instances ei
		WHERE ei.id = $1
	`

	var exerciseInstance exerciseInstanceEntity
	err := pgxscan.Get(ctx, r.pool, &exerciseInstance, query, uuidToPgtype(id))
	if err != nil {
		logger.Errorf("failed to get exercise instance by id: %v", err)
		if err == pgx.ErrNoRows {
			return domain.ExerciseInstance{}, domain.ErrNotFound
		}
		return domain.ExerciseInstance{}, err
	}

	return exerciseInstance.toDomain(), nil
}

func (r *PGXRepository) GetExerciseInstancesByRoutineID(ctx context.Context, routineID domain.ID) ([]domain.ExerciseInstance, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetExerciseInstancesByRoutineID")
	defer span.Finish()

	query := `
		SELECT * FROM exercise_instances ei
		WHERE ei.routine_id = $1
	`

	var exerciseInstances []exerciseInstanceEntity
	err := pgxscan.Select(ctx, r.pool, &exerciseInstances, query, uuidToPgtype(routineID))
	if err != nil {
		logger.Errorf("failed to get exercise instances by routine id: %v", err)
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	result := make([]domain.ExerciseInstance, len(exerciseInstances))
	for i, exerciseInstance := range exerciseInstances {
		result[i] = exerciseInstance.toDomain()
	}

	return result, nil
}

func (r *PGXRepository) CreateExerciseInstance(ctx context.Context, exerciseInstance domain.ExerciseInstance) (domain.ExerciseInstance, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.CreateExerciseInstance")
	defer span.Finish()

	query := `
		INSERT INTO exercise_instances (id, routine_id, exercise_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`

	entity := exerciseInstanceFromDomain(exerciseInstance)

	err := pgxscan.Get(ctx, r.pool, &entity, query, entity.ID, entity.RoutineID, entity.ExerciseID, entity.CreatedAt)
	if err != nil {
		logger.Errorf("failed to create exercise instance: %v", err)
		return domain.ExerciseInstance{}, err
	}

	return entity.toDomain(), nil
}

func (r *PGXRepository) DeleteExerciseInstance(ctx context.Context, id domain.ID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.DeleteExerciseInstance")
	defer span.Finish()

	query := `
		DELETE FROM exercise_instances
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, uuidToPgtype(id))
	if err != nil {
		logger.Errorf("failed to delete exercise instance: %v", err)
		return err
	}

	return nil
}
