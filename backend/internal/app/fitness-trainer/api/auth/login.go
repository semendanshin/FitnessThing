package auth

import (
	"context"
	"fmt"

	"fitness-trainer/internal/app/mappers"
	"fitness-trainer/internal/domain"
	desc "fitness-trainer/pkg/workouts"

	"github.com/opentracing/opentracing-go"
)

func (i *Implementation) Login(ctx context.Context, in *desc.LoginRequest) (*desc.LoginResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.auth.Login")
	defer span.Finish()

	if err := in.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInvalidArgument, err)
	}

	tokens, err := i.service.Login(ctx, in.Email, in.Password)
	if err != nil {
		return nil, err
	}

	return &desc.LoginResponse{
		Tokens: mappers.TokensToProto(tokens),
	}, nil
}
