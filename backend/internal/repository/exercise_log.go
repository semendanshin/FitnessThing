package repository

import (
	"context"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/logger"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/opentracing/opentracing-go"
)

type exerciseLogEntity struct {
	ID          pgtype.UUID        `db:"id"`
	ExerciseID  pgtype.UUID        `db:"exercise_id"`
	WorkoutID   pgtype.UUID        `db:"workout_id"`
	Notes       string             `db:"notes"`
	PowerRating int                `db:"power_rating"`
	CreatedAt   pgtype.Timestamptz `db:"created_at"`
	UpdateAt    pgtype.Timestamptz `db:"updated_at"`
}

func (e exerciseLogEntity) toDomain() domain.ExerciseLog {
	return domain.ExerciseLog{
		Model: domain.Model{
			ID:        domain.ID(e.ID.Bytes),
			CreatedAt: e.CreatedAt.Time,
			UpdatedAt: e.UpdateAt.Time,
		},
		ExerciseID:  domain.ID(e.ExerciseID.Bytes),
		WorkoutID:   domain.ID(e.WorkoutID.Bytes),
		Notes:       e.Notes,
		PowerRating: e.PowerRating,
	}
}

func exerciseLogFromDomain(exerciseLog domain.ExerciseLog) exerciseLogEntity {
	return exerciseLogEntity{
		ID:          uuidToPgtype(exerciseLog.ID),
		ExerciseID:  uuidToPgtype(exerciseLog.ExerciseID),
		WorkoutID:   uuidToPgtype(exerciseLog.WorkoutID),
		Notes:       exerciseLog.Notes,
		PowerRating: exerciseLog.PowerRating,
		CreatedAt:   timeToPgtype(exerciseLog.CreatedAt),
		UpdateAt:    timeToPgtype(exerciseLog.UpdatedAt),
	}
}

func toExerciseLogsDomain(exerciseLogs []exerciseLogEntity) []domain.ExerciseLog {
	var result []domain.ExerciseLog
	for _, exerciseLog := range exerciseLogs {
		result = append(result, exerciseLog.toDomain())
	}

	return result
}

func (r *PGXRepository) GetExerciseLogsByWorkoutID(ctx context.Context, workoutID domain.ID) ([]domain.ExerciseLog, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetExerciseLogsByWorkoutID")
	defer span.Finish()

	query := `
		SELECT id, exercise_id, workout_id, notes, power_rating, created_at, updated_at
		FROM exercise_logs
		WHERE workout_id = $1
		ORDER BY created_at
	`

	engine := r.contextManager.GetEngineFromContext(ctx)

	var exerciseLogs []exerciseLogEntity
	if err := pgxscan.Select(ctx, engine, &exerciseLogs, query, uuidToPgtype(workoutID)); err != nil {
		return nil, err
	}

	return toExerciseLogsDomain(exerciseLogs), nil
}

func (r *PGXRepository) CreateExerciseLog(ctx context.Context, exerciseLog domain.ExerciseLog) (domain.ExerciseLog, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.CreateExerciseLog")
	defer span.Finish()

	query := `
		INSERT INTO exercise_logs (id, exercise_id, workout_id, notes, power_rating, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *
	`

	engine := r.contextManager.GetEngineFromContext(ctx)

	exerciseLogEntity := exerciseLogFromDomain(exerciseLog)
	if err := pgxscan.Get(ctx, engine, &exerciseLogEntity, query, exerciseLogEntity.ID, exerciseLogEntity.ExerciseID, exerciseLogEntity.WorkoutID, exerciseLogEntity.Notes, exerciseLogEntity.PowerRating, exerciseLogEntity.CreatedAt); err != nil {
		return domain.ExerciseLog{}, err
	}

	return exerciseLogEntity.toDomain(), nil
}

func (r *PGXRepository) GetExerciseLogByID(ctx context.Context, id domain.ID) (domain.ExerciseLog, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetExerciseLogByID")
	defer span.Finish()

	query := `
		SELECT id, exercise_id, workout_id, notes, power_rating, created_at, updated_at
		FROM exercise_logs
		WHERE id = $1
	`

	engine := r.contextManager.GetEngineFromContext(ctx)

	var exerciseLog exerciseLogEntity
	if err := pgxscan.Get(ctx, engine, &exerciseLog, query, uuidToPgtype(id)); err != nil {
		return domain.ExerciseLog{}, err
	}

	return exerciseLog.toDomain(), nil
}

func (r *PGXRepository) GetExerciseLogsByExerciseIDAndUserID(ctx context.Context, exerciseID, userID domain.ID) ([]domain.ExerciseLog, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetExerciseLogsByExerciseIDAndUserID")
	defer span.Finish()

	query := `
		SELECT el.id, el.exercise_id, el.workout_id, el.notes, el.power_rating, el.created_at, el.updated_at
		FROM exercise_logs el
		JOIN workouts w ON el.workout_id = w.id
		WHERE el.exercise_id = $1 AND w.user_id = $2
		ORDER BY el.created_at DESC;
	`

	engine := r.contextManager.GetEngineFromContext(ctx)

	var exerciseLogs []domain.ExerciseLog
	err := pgxscan.Select(ctx, engine, &exerciseLogs, query, exerciseID, userID)
	if err != nil {
		logger.Errorf("failed to get exercise logs by exercise id and user id: %v", err)
		return nil, err
	}

	return exerciseLogs, nil
}

func (r *PGXRepository) DeleteExerciseLog(ctx context.Context, id domain.ID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.DeleteExerciseLog")
	defer span.Finish()

	query := `
		DELETE FROM exercise_logs
		WHERE id = $1
	`

	engine := r.contextManager.GetEngineFromContext(ctx)

	_, err := engine.Exec(ctx, query, uuidToPgtype(id))
	if err != nil {
		return err
	}

	return nil
}

func (r *PGXRepository) UpdateExerciseLog(ctx context.Context, id domain.ID, exerciseLog domain.ExerciseLog) (domain.ExerciseLog, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.UpdateExerciseLog")
	defer span.Finish()

	query := `
		UPDATE exercise_logs
		SET notes = $1, power_rating = $2, updated_at = now()
		WHERE id = $3
		RETURNING *
	`

	engine := r.contextManager.GetEngineFromContext(ctx)

	exerciseLogEntity := exerciseLogFromDomain(exerciseLog)
	if err := pgxscan.Get(ctx, engine, &exerciseLogEntity, query, exerciseLogEntity.Notes, exerciseLogEntity.PowerRating, exerciseLogEntity.ID); err != nil {
		return domain.ExerciseLog{}, err
	}

	return exerciseLogEntity.toDomain(), nil
}
