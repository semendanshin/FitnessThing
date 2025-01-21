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

type exerciseEntity struct {
	ID                 pgtype.UUID
	Name               string
	Description        pgtype.Text
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

	var exercises []exerciseEntity
	err := pgxscan.Select(ctx, r.pool, &exercises, query, muscleGroups, excludedExercises)
	if err != nil {
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

	var exercise exerciseEntity
	err := pgxscan.Get(ctx, r.pool, &exercise, query, id)
	if err != nil {
		logger.Errorf("failed to get exercise by id: %v", err)
		return domain.Exercise{}, err
	}

	return exercise.toDomain(), nil
}

func (r *PGXRepository) CreateExercise(ctx context.Context, exercise domain.Exercise) (domain.Exercise, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.CreateExercise")
	defer span.Finish()
	
	exerciseQuery := `
		INSERT INTO exercises (id, name, description, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`

	exercieMuscleGroupQuery := `
		INSERT INTO exercise_muscle_groups (exercise_id, muscle_group_id)
		WITH muscle_groups AS (
			SELECT id
			FROM muscle_groups
			WHERE name = ANY($1)
		)
		SELECT $2, id
		FROM muscle_groups
		RETURNING exercise_id, muscle_group_id
	`

	err := r.runInTransaction(ctx, func(tx pgx.Tx) error {
		exerciseEntity := exerciseFromDomain(exercise)

		err := pgxscan.Get(
			ctx,
			r.pool,
			&exerciseEntity,
			exerciseQuery,
			exerciseEntity.ID,
			exerciseEntity.Name,
			exerciseEntity.Description,
			exerciseEntity.CreatedAt,
		)
		if err != nil {
			logger.Errorf("failed to create exercise: %v", err)
			return err
		}

		type exerciseMuscleGroup struct {
			ExerciseID   pgtype.UUID
			MuscleGroupID pgtype.UUID
		}

		for _, mg := range exercise.TargetMuscleGroups {
			var emg exerciseMuscleGroup
			err := pgxscan.Get(ctx, tx, &emg, exercieMuscleGroupQuery, mg, exerciseEntity.ID)
			if err != nil {
				logger.Errorf("failed to create exercise muscle group: %v", err)
				return err
			}
		}

		return nil
	})

	if err != nil {
		return domain.Exercise{}, err
	}

	return exercise, nil
}
