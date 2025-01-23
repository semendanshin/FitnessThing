package interceptors

import (
	"context"
	"fitness-trainer/internal/logger"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
)

func RecovertInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("[interceptor.Recovery] method: %s; error: %v", info.FullMethod, r)
			span := opentracing.SpanFromContext(ctx)
			if span == nil {
				return
			}
			ext.Error.Set(span, true)
			span.SetTag("error.message", r)
		}
	}()

	return handler(ctx, req)
}
