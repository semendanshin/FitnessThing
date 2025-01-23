package dto

import (
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/utils"
)

type CreateExerciseDTO struct {
	Name               string
	Description        utils.Nullable[string]
	VideoURL           utils.Nullable[string]
	TargetMuscleGroups []domain.ID
}
