package auth

import (
	"context"
	"fmt"

	"fitness-trainer/internal/domain"
	desc "fitness-trainer/pkg/workouts"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opentracing/opentracing-go"
)

func (i *Implementation) Logout(ctx context.Context, in *desc.LogoutRequest) (*empty.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.auth.logout")
	defer span.Finish()

	if err := in.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInvalidArgument, err)
	}

	err := i.service.Logout(ctx, in.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
