package dto

import "fitness-trainer/internal/domain"

type GenerateWorkoutOptions struct {
	UserID         domain.ID
	Workouts       []SlimWorkoutDTO
	Exercises      []SlimExerciseDTO
	VarietyLevel   int
	UserPrompt     string
	BaseUserPrompt string
}
