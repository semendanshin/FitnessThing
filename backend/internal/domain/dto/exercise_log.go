package dto

import "fitness-trainer/internal/domain"

type ExerciseLogDTO struct {
	ExerviceLog domain.ExerciseLog
	Exercise    domain.Exercise
	SetLogs     []domain.ExerciseSetLog
}
