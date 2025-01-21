package user

import (
	"context"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"
	desc "fitness-trainer/pkg/workouts"
)

type Service interface {
	CreateUser(ctx context.Context, dto dto.CreateUserDTO) (domain.User, error)
	GetUserByID(ctx context.Context, id domain.ID) (domain.User, error)
	UpdateUser(ctx context.Context, id domain.ID, dto dto.UpdateUserDTO) (domain.User, error)
}

type Implementation struct {
	service Service
	desc.UnimplementedUserServiceServer
}

func New(service Service) *Implementation {
	return &Implementation{
		service: service,
	}
}
