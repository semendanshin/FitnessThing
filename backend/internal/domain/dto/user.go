package dto

import (
	"fitness-trainer/internal/utils"
	"time"
)

type CreateUserDTO struct {
	Email    string
	Password string

	DateOfBirth time.Time
	FirstName   utils.Nullable[string]
	LastName    utils.Nullable[string]
	Height      utils.Nullable[float32]
	Weight      utils.Nullable[float32]
}

type UpdateUserDTO struct {
	FirstName     utils.Nullable[string]
	LastName      utils.Nullable[string]
	Height        utils.Nullable[float32]
	Weight        utils.Nullable[float32]
	ProfilePicURL utils.Nullable[string]
	DateOfBirth   time.Time
}
