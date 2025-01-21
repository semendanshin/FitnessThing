package service

import (
	"context"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"

	"github.com/opentracing/opentracing-go"
)

func (s *Service) GetExercises(ctx context.Context, muscleGroups, excludedExercises []domain.ID) ([]domain.Exercise, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetExercises")
	defer span.Finish()

	return s.exerciseRepository.GetExercises(ctx, muscleGroups, excludedExercises)
}

func (s *Service) GetExerciseByID(ctx context.Context, id domain.ID) (domain.Exercise, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetExerciseByID")
	defer span.Finish()

	return s.exerciseRepository.GetExerciseByID(ctx, id)
}

func (s *Service) GetExerciseAlternatives(ctx context.Context, id domain.ID) ([]domain.Exercise, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetExerciseAlternatives")
	defer span.Finish()

	exercise, err := s.GetExerciseByID(ctx, id)
	if err != nil {
		return nil, err
	}

	ids := make([]domain.ID, 0, len(exercise.TargetMuscleGroups))
	for _, muscleGroup := range exercise.TargetMuscleGroups {
		mg, err := s.muscleGroupRepository.GetMuscleGroupByName(ctx, muscleGroup.String())
		if err != nil {
			return nil, err
		}
		ids = append(ids, mg.ID)
	}

	result, err := s.exerciseRepository.GetExercises(ctx, ids, []domain.ID{id})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) GetExerciseHistory(ctx context.Context, userID, exerciseID domain.ID) ([]dto.ExerciseLogDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetExerciseHistory")
	defer span.Finish()

	exerciseLogs, err := s.exerciseLogRepository.GetExerciseLogsByExerciseIDAndUserID(ctx, exerciseID, userID)
	if err != nil {
		return nil, err
	}

	exerciseLogDTOs := make([]dto.ExerciseLogDTO, 0, len(exerciseLogs))
	for _, exerciseLog := range exerciseLogs {
		exerciseLogDTO, err := s.GetExerciseLog(ctx, userID, exerciseLog.ID)
		if err != nil {
			return nil, err
		}

		exerciseLogDTOs = append(exerciseLogDTOs, exerciseLogDTO)
	}

	return exerciseLogDTOs, nil
}
