package interceptors

import (
	"context"
	"errors"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrCodesInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "%s", err.Error())
		}
		if errors.Is(err, domain.ErrAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
		}
		if errors.Is(err, domain.ErrInvalidArgument) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorized) {
			return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
		}
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Errorf(codes.PermissionDenied, "%s", err.Error())
		}

		logger.Errorf("[interceptor.Error] method: %s; error: %s", info.FullMethod, err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return resp, err
}
