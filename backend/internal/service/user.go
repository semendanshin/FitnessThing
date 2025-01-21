package service

import (
	"context"
	"time"

	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"
	"fitness-trainer/internal/utils"

	"github.com/opentracing/opentracing-go"
)

func (s *Service) CreateUser(ctx context.Context, dto dto.CreateUserDTO) (domain.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.CreateUser")
	defer span.Finish()

	hashedPass, err := utils.HashPassword(dto.Password)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.NewUser(
		dto.Email,
		hashedPass,
		dto.FirstName.V,
		dto.LastName.V,
		dto.DateOfBirth,
		dto.Height.V,
		dto.Weight.V,
	)

	user, err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, id domain.ID) (domain.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetUserByID")
	defer span.Finish()

	return s.userRepository.GetUserByID(ctx, id)
}

func (s *Service) UpdateUser(ctx context.Context, id domain.ID, dto dto.UpdateUserDTO) (domain.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.UpdateUser")
	defer span.Finish()
	
	user, err := s.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	{
		if !dto.DateOfBirth.IsZero() {
			user.DateOfBirth = dto.DateOfBirth
		}

		if dto.LastName.IsValid {
			user.LastName = dto.LastName.V
		}

		if dto.FirstName.IsValid {
			user.FirstName = dto.FirstName.V
		}

		if dto.Height.IsValid {
			user.Height = dto.Height.V
		}

		if dto.Weight.IsValid {
			user.Weight = dto.Weight.V
		}

		user.UpdatedAt = time.Now()
	}

	return s.userRepository.UpdateUser(ctx, user)
}
