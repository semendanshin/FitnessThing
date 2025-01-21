package dto

import "fitness-trainer/internal/domain"

type WorkoutDetailsDTO struct {
	Workout      domain.Workout
	ExerciseLogs []ExerciseLogDTO
}
