package mappers

import (
	"fitness-trainer/internal/domain"
	desc "fitness-trainer/pkg/workouts"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserToProto(user domain.User) *desc.User {
	return &desc.User{
		Id:          user.ID.String(),
		Email:       user.Email,
		DateOfBirth: timestamppb.New(user.DateOfBirth),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Weight:      user.Weight,
		Height:      user.Height,
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
	}
}
