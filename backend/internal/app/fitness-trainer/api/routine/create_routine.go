package routine

import (
	"context"
	"fmt"

	"fitness-trainer/internal/app/interceptors"
	"fitness-trainer/internal/app/mappers"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"
	"fitness-trainer/internal/logger"
	desc "fitness-trainer/pkg/workouts"

	"github.com/opentracing/opentracing-go"
)

func (i *Implementation) CreateRoutine(ctx context.Context, in *desc.CreateRoutineRequest) (*desc.RoutineResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.routine.CreateRoutine")
	defer span.Finish()

	if err := in.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInvalidArgument, err)
	}

	var dto dto.CreateRoutineDTO
	{
		var ok bool
		dto.UserID, ok = interceptors.GetUserID(ctx)
		if !ok {
			logger.Errorf("user id is not found in context")
			return nil, domain.ErrUnauthorized
		}

		dto.Name = in.GetName()
		dto.Description = in.GetDescription()
	}

	routine, err := i.service.CreateRoutine(ctx, dto)
	if err != nil {
		logger.Errorf("error creating routine: %v", err)
		return nil, err
	}

	return &desc.RoutineResponse{
		Routine: mappers.RoutineToProto(routine),
	}, nil
}
