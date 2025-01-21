package workout

import (
	"context"

	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"

	desc "fitness-trainer/pkg/workouts"
)

type Service interface {
	StartWorkout(ctx context.Context, userID domain.ID, routineID *domain.ID) (domain.Workout, error)
	GetWorkout(ctx context.Context, userID domain.ID, workoutID domain.ID) (dto.WorkoutDetailsDTO, error)
	LogExercise(ctx context.Context, userID, workoutID, exerciseID domain.ID) (domain.ExerciseLog, error)
	GetExerciseLog(ctx context.Context, userID, exerciseLogID domain.ID) (dto.ExerciseLogDTO, error)
	LogSet(ctx context.Context, userID, workoutID, exerciseLogID domain.ID, setlogDTO dto.CreateSetLogDTO) (domain.ExerciseSetLog, error)
	GetActiveWorkouts(ctx context.Context, userID domain.ID) ([]domain.Workout, error)
	CompleteWorkout(ctx context.Context, userID domain.ID, workoutID domain.ID) error
}

type Implementation struct {
	service Service
	desc.UnimplementedWorkoutServiceServer
}

func New(service Service) *Implementation {
	return &Implementation{
		service: service,
	}
}
