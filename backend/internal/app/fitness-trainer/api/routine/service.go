package routine

import (
	"context"

	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"
	desc "fitness-trainer/pkg/workouts"
)

type Service interface {
	GetRoutines(ctx context.Context, userID domain.ID) ([]domain.Routine, error)
	CreateRoutine(ctx context.Context, dto dto.CreateRoutineDTO) (domain.Routine, error)
	GetRoutineByID(ctx context.Context, id domain.ID) (dto.RoutineDetailsDTO, error)
	DeleteRoutine(ctx context.Context, id domain.ID) error

	AddExerciseToRoutine(ctx context.Context, routineID domain.ID, exerciseID domain.ID) (domain.ExerciseInstance, error)
	RemoveExerciseInstanceFromRoutine(ctx context.Context, userID, routineID, exerciseInstanceID domain.ID) error
}

type Implementation struct {
	service Service
	desc.UnimplementedRoutineServiceServer
}

func New(service Service) *Implementation {
	return &Implementation{
		service: service,
	}
}
