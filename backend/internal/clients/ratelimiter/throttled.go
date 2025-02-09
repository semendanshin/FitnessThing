package ratelimiter

import (
	"context"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/logger"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/throttled/throttled/v2"
)

type RateLimiter struct {
	rateLimiter *throttled.GCRARateLimiterCtx
}

func New(rateLimiter *throttled.GCRARateLimiterCtx) *RateLimiter {
	return &RateLimiter{
		rateLimiter: rateLimiter,
	}
}

func (r *RateLimiter) Allow(ctx context.Context, userID domain.ID) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ratelimiter.Allow")
	defer span.Finish()

	exceeded, result, err := r.rateLimiter.RateLimitCtx(ctx, userID.String(), 1)
	if err != nil {
		return false, fmt.Errorf("failed to check rate limit: %w", err)
	}

	logger.Debugf("rate limit result: %+v", result)
	logger.Debugf("rate limit allowed: %v", exceeded)

	return !exceeded, nil
}
