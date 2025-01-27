package workout

import (
	"context"

	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"

	desc "fitness-trainer/pkg/workouts"
)

type Service interface {
	GetWorkouts(ctx context.Context, userID domain.ID, limit, offset int) ([]dto.WorkoutDTO, error)
	StartWorkout(ctx context.Context, userID domain.ID, routineID *domain.ID) (domain.Workout, error)
	GetWorkout(ctx context.Context, userID, workoutID domain.ID) (dto.WorkoutDetailsDTO, error)
	DeleteWorkout(ctx context.Context, userID, workoutID domain.ID) error
	GetActiveWorkouts(ctx context.Context, userID domain.ID) ([]domain.Workout, error)
	CompleteWorkout(ctx context.Context, userID, workoutID domain.ID) error
	RateWorkout(ctx context.Context, userID, workoutID domain.ID, rating int) (domain.Workout, error)
	AddCommentToWorkout(ctx context.Context, userID, workoutID domain.ID, comment string) (domain.Workout, error)

	LogExercise(ctx context.Context, userID, workoutID, exerciseID domain.ID) (domain.ExerciseLog, error)
	GetExerciseLog(ctx context.Context, userID, exerciseLogID domain.ID) (dto.ExerciseLogDTO, error)
	DeleteExerciseLog(ctx context.Context, userID, workoutID, exerciseLogID domain.ID) error

	LogSet(ctx context.Context, userID, workoutID, exerciseLogID domain.ID, setlogDTO dto.CreateSetLogDTO) (domain.ExerciseSetLog, error)
	DeleteSetLog(ctx context.Context, userID, workoutID, exerciseLogID domain.ID, setLogID domain.ID) error
	UpdateSetLog(ctx context.Context, userID, workoutID, exerciseLogID, setLogID domain.ID, setlogDTO dto.UpdateSetLogDTO) (domain.ExerciseSetLog, error)
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
