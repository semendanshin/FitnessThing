package interceptors

import (
	"context"
	"fitness-trainer/internal/logger"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func RecovertInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "interceptor.Recovery")
	defer span.Finish()

	defer func() {
		if r := recover(); r != nil {
			span.SetTag("error", true)
			span.SetTag("error.message", r)
			logger.Errorf("[interceptor.Recovery] method: %s; error: %v", info.FullMethod, r)
		}
	}()

	return handler(ctx, req)
}
