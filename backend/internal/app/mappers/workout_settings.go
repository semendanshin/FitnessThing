package mappers

import (
	"fitness-trainer/internal/domain"
	desc "fitness-trainer/pkg/workouts"
)

func GenerationSettingsToProto(settings domain.GenerationSettings) *desc.WorkoutGenerationSettings {
	return &desc.WorkoutGenerationSettings{
		BasePrompt:   settings.BasePrompt,
		VarietyLevel: int32(settings.VarietyLevel),
	}
}
