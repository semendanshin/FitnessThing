package user

import (
	"context"
	"fmt"

	"fitness-trainer/internal/app/mappers"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"
	"fitness-trainer/internal/logger"
	"fitness-trainer/internal/utils"

	desc "fitness-trainer/pkg/workouts"

	"github.com/opentracing/opentracing-go"
)

func (i *Implementation) CreateUser(ctx context.Context, in *desc.CreateUserRequest) (*desc.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.user.create")
	defer span.Finish()

	if err := in.Validate(); err != nil {
		logger.Errorf("error validating request: %v", err)
		return nil, fmt.Errorf("%w: %w", domain.ErrInvalidArgument, err)
	}

	var input dto.CreateUserDTO
	{
		input.Email = in.Email
		input.Password = in.Password

		input.FirstName = utils.NewNullable(in.GetFirstName(), in.GetFirstName() != "")
		input.LastName = utils.NewNullable(in.GetLastName(), in.GetLastName() != "")

		input.DateOfBirth = in.GetDateOfBirth().AsTime()

		input.Height = utils.NewNullable(in.GetHeight(), in.GetHeight() != 0)
		input.Weight = utils.NewNullable(in.GetWeight(), in.GetWeight() != 0)
	}

	user, err := i.service.CreateUser(ctx, input)
	if err != nil {
		logger.Errorf("error creating user: %v", err)
		return nil, err
	}

	return &desc.UserResponse{
		User: mappers.UserToProto(user),
	}, nil
}
