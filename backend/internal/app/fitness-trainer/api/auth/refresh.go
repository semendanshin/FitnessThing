package auth

import (
	"context"
	"fmt"

	"fitness-trainer/internal/app/mappers"
	"fitness-trainer/internal/domain"
	desc "fitness-trainer/pkg/workouts"

	"github.com/opentracing/opentracing-go"
)

func (i *Implementation) Refresh(ctx context.Context, in *desc.RefreshRequest) (*desc.RefreshResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.auth.refresh")
	defer span.Finish()

	if err := in.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInvalidArgument, err)
	}

	tokens, err := i.service.Refresh(ctx, mappers.ProtoToTokens(in.Tokens))
	if err != nil {
		return nil, err
	}

	return &desc.RefreshResponse{
		Tokens: mappers.TokensToProto(tokens),
	}, nil
}
