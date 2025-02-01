package dto

import "fitness-trainer/internal/domain"

type ExerciseLogDTO struct {
	ExerciseLog  domain.ExerciseLog
	Exercise     domain.Exercise
	SetLogs      []domain.ExerciseSetLog
	ExpectedSets []domain.ExpectedSet
}
