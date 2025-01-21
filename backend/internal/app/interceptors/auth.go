package interceptors

import (
	"context"

	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/logger"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	userIDKey            contextKey = "auth-interceptor.user-id"
	accesTokenHeaderName string     = "x-access-token"
)

type JWTManager interface {
	ParseToken(ctx context.Context, token string) (domain.ID, error)
}

func NewAuth(
	jwtManager JWTManager,
	unprotectedMethods map[string]struct{},
) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		span, ctx := opentracing.StartSpanFromContext(ctx, "interceptors.Auth")
		defer span.Finish()

		logger.Debugf("method %s is called", info.FullMethod)
		if _, ok := unprotectedMethods[info.FullMethod]; ok {
			logger.Infof("method %s is unprotected", info.FullMethod)
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Errorf("metadata is not provided")
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		token, ok := md[accesTokenHeaderName]
		if !ok {
			logger.Errorf("authorization token is not provided")
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		id, err := jwtManager.ParseToken(ctx, token[0])
		if err != nil {
			logger.Errorf("invalid token: %v", err)
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, userIDKey, id)

		return handler(ctx, req)
	}
}

func GetUserID(ctx context.Context) (domain.ID, bool) {
	value, ok := ctx.Value(userIDKey).(domain.ID)
	return value, ok
}
