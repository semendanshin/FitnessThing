package exercise

import (
	"context"

	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"
	desc "fitness-trainer/pkg/workouts"
)

type Service interface {
	CreateExercise(ctx context.Context, exercise dto.CreateExerciseDTO) (domain.Exercise, error)
	GetExercises(ctx context.Context, muscleGroups, excludedExercises []domain.ID) ([]domain.Exercise, error)
	GetExerciseByID(ctx context.Context, id domain.ID) (domain.Exercise, error)
	GetExerciseAlternatives(ctx context.Context, id domain.ID) ([]domain.Exercise, error)

	GetMuscleGroups(ctx context.Context) ([]dto.MuscleGroupDTO, error)

	GetExerciseHistory(ctx context.Context, userID, exerciseID domain.ID) ([]dto.ExerciseLogDTO, error)
}

type Implementation struct {
	service Service
	desc.UnimplementedExerciseServiceServer
}

func New(service Service) *Implementation {
	return &Implementation{
		service: service,
	}
}
