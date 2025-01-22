package dto

import (
	"time"

	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/utils"
)

type RoutineDetailsDTO struct {
	ID                domain.ID
	UserID            domain.ID
	Name              string
	Description       string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ExerciseInstances []ExerciseInstanceDetailsDTO
}

type CreateRoutineDTO struct {
	UserID      domain.ID
	WorkoutID   utils.Nullable[domain.ID]
	Name        string
	Description string
}

type UpdateRoutineDTO struct {
	Name        utils.Nullable[string]
	Description utils.Nullable[string]
}
