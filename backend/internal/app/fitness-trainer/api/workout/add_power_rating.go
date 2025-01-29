package workout

import (
	"context"
	"fitness-trainer/internal/app/interceptors"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/logger"
	desc "fitness-trainer/pkg/workouts"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) AddPowerRatingToExerciseLog(ctx context.Context, in *desc.AddPowerRatingToExerciseLogRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.workout.AddPowerRating")
	defer span.Finish()

	if err := in.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInvalidArgument, err)
	}

	userID, ok := interceptors.GetUserID(ctx)
	if !ok {
		logger.Errorf("error getting user id from context")
		return nil, domain.ErrInternal
	}

	workoutID, err := domain.ParseID(in.GetWorkoutId())
	if err != nil {
		logger.Errorf("error parsing workout id %s: %v", in.GetWorkoutId(), err)
		return nil, fmt.Errorf("%w: %w", domain.ErrInvalidArgument, err)
	}

	exerciseLogID, err := domain.ParseID(in.GetExerciseLogId())
	if err != nil {
		logger.Errorf("error parsing exercise log id %s: %v", in.GetExerciseLogId(), err)
		return nil, fmt.Errorf("%w: %w", domain.ErrInvalidArgument, err)
	}

	err = i.service.AddPowerRatingToExerciseLog(ctx, userID, workoutID, exerciseLogID, int(in.GetPowerRating()))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
