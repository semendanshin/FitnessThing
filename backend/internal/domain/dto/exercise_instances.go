package dto

import (
	"fitness-trainer/internal/domain"
	"time"
)

type ExerciseInstanceDetailsDTO struct {
	ID         domain.ID
	RoutineID  domain.ID
	ExerciseID domain.ID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Exercise   domain.Exercise
	Sets       []domain.Set
}
