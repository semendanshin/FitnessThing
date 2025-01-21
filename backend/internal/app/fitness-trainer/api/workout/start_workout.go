package workout

import (
	"context"
	"fitness-trainer/internal/app/interceptors"
	"fitness-trainer/internal/app/mappers"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/logger"
	desc "fitness-trainer/pkg/workouts"
	"fmt"

	"github.com/opentracing/opentracing-go"
)

func (i *Implementation) StartWorkout(ctx context.Context, in *desc.StartWorkoutRequest) (*desc.WorkoutResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.workout.StartWorkout")
	defer span.Finish()

	if err := in.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInvalidArgument, err)
	}

	userID, ok := interceptors.GetUserID(ctx)
	if !ok {
		logger.Errorf("user id not found in context")
		return nil, domain.ErrInternal
	}

	var routineID *domain.ID
	var err error
	if in.RoutineId != nil {
		parsedID, err := domain.ParseID(*in.RoutineId)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", domain.ErrInvalidArgument, err)
		}
		routineID = &parsedID
	}

	workout, err := i.service.StartWorkout(ctx, userID, routineID)
	if err != nil {
		return nil, err
	}

	return &desc.WorkoutResponse{
		Workout: mappers.WorkoutToProto(workout),
	}, nil
}
