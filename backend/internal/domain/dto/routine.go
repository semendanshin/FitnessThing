package dto

import (
	"time"
	
	"fitness-trainer/internal/domain"
)

type RoutineDetailsDTO struct {
	ID          domain.ID
	UserID      domain.ID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExerciseInstances []ExerciseInstanceDetailsDTO
}

type CreateRoutineDTO struct {
	UserID      domain.ID
	Name        string
	Description string
}
