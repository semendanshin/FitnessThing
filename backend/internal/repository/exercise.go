package repository

import (
	"context"
	"errors"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/logger"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/opentracing/opentracing-go"
)

type exerciseEntity struct {
	ID                 pgtype.UUID
	Name               string
	Description        pgtype.Text
	VideoURL           pgtype.Text
	TargetMuscleGroups pgtype.Array[string]
	CreatedAt          pgtype.Timestamptz
	UpdatedAt          pgtype.Timestamptz
}

func (e exerciseEntity) toDomain() domain.Exercise {
	musclegroups := make([]domain.MuscleGroup, len(e.TargetMuscleGroups.Elements))
	for i, mg := range e.TargetMuscleGroups.Elements {
		musclegroups[i] = domain.MuscleGroup(mg)
	}
	return domain.Exercise{
		Model: domain.Model{
			ID:        domain.ID(e.ID.Bytes),
			CreatedAt: e.CreatedAt.Time,
			UpdatedAt: e.UpdatedAt.Time,
		},
		Name:               e.Name,
		Description:        e.Description.String,
		VideoURL:           e.VideoURL.String,
		TargetMuscleGroups: musclegroups,
	}
}

func exerciseFromDomain(exercise domain.Exercise) exerciseEntity {
	musclegroups := make([]string, len(exercise.TargetMuscleGroups))
	for i, mg := range exercise.TargetMuscleGroups {
		musclegroups[i] = string(mg)
	}
	return exerciseEntity{
		ID:                 uuidToPgtype(exercise.ID),
		Name:               exercise.Name,
		Description:        pgtype.Text{String: exercise.Description, Valid: exercise.Description != ""},
		VideoURL:           pgtype.Text{String: exercise.VideoURL, Valid: exercise.VideoURL != ""},
		TargetMuscleGroups: pgtype.Array[string]{Elements: musclegroups, Valid: true},
		CreatedAt:          timeToPgtype(exercise.CreatedAt),
		UpdatedAt:          timeToPgtype(exercise.UpdatedAt),
	}
}

func (r *PGXRepository) GetExercises(ctx context.Context, muscleGroups, excludedExercises []domain.ID) ([]domain.Exercise, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetExercises")
	defer span.Finish()

	query := `
		SELECT e.id, e.name, e.description, e.created_at, ARRAY_AGG(mg.name) AS target_muscle_groups, e.updated_at
		FROM exercise_muscle_groups emg
		JOIN exercises e ON emg.exercise_id = e.id
		JOIN muscle_groups mg ON emg.muscle_group_id = mg.id
		WHERE e.id NOT IN (SELECT UNNEST($2::UUID[]))
		GROUP BY e.id
		HAVING ARRAY_AGG(mg.id) && $1 OR $1 = '{}'
		ORDER BY e.created_at DESC;
	`

	engine := r.contextManager.GetEngineFromContext(ctx)

	var exercises []exerciseEntity
	err := pgxscan.Select(ctx, engine, &exercises, query, muscleGroups, excludedExercises)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.Exercise{}, nil
		}
		logger.Errorf("failed to get exercises: %v", err)
		return nil, err
	}

	result := make([]domain.Exercise, len(exercises))
	for i, e := range exercises {
		result[i] = e.toDomain()
	}

	return result, nil
}

func (r *PGXRepository) GetExerciseByID(ctx context.Context, id domain.ID) (domain.Exercise, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetExerciseByID")
	defer span.Finish()

	query := `
		SELECT e.id, e.name, e.description, e.created_at, ARRAY_AGG(mg.name) AS target_muscle_groups, e.updated_at
		FROM exercises e
		JOIN exercise_muscle_groups emg ON e.id = emg.exercise_id
		JOIN muscle_groups mg ON emg.muscle_group_id = mg.id
		WHERE e.id = $1
		GROUP BY e.id;
	`

	engine := r.contextManager.GetEngineFromContext(ctx)

	var exercise exerciseEntity
	err := pgxscan.Get(ctx, engine, &exercise, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Exercise{}, domain.ErrNotFound
		}
		logger.Errorf("failed to get exercise by id: %v", err)
		return domain.Exercise{}, err
	}

	return exercise.toDomain(), nil
}

// CreateExercise must be called within a transaction
func (r *PGXRepository) CreateExercise(ctx context.Context, exercise domain.Exercise, muscleGroupIDs []domain.ID) (domain.Exercise, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.CreateExercise")
	defer span.Finish()

	exerciseQuery := `
		INSERT INTO exercises (id, name, description, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`

	exercieMuscleGroupsQuery := `
		INSERT INTO exercise_muscle_groups (muscle_group_id, exercise_id)
		SELECT UNNEST($1::UUID[]), $2
		RETURNING *
	`

	convertedMuscleGroupIDs := make([]pgtype.UUID, len(muscleGroupIDs))
	for i, id := range muscleGroupIDs {
		convertedMuscleGroupIDs[i] = uuidToPgtype(id)
	}

	engine := r.contextManager.GetEngineFromContext(ctx)

	exerciseEntity := exerciseFromDomain(exercise)

	err := pgxscan.Get(
		ctx,
		engine,
		&exerciseEntity,
		exerciseQuery,
		exerciseEntity.ID,
		exerciseEntity.Name,
		exerciseEntity.Description,
		exerciseEntity.CreatedAt,
	)
	if err != nil {
		return domain.Exercise{}, nil
	}

	_, err = engine.Exec(
		ctx,
		exercieMuscleGroupsQuery,
		convertedMuscleGroupIDs,
		exerciseEntity.ID,
	)
	if err != nil {
		return domain.Exercise{}, nil
	}

	return r.GetExerciseByID(ctx, exercise.ID)
}
