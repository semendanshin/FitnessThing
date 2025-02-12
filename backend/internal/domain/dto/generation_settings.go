package dto

import "fitness-trainer/internal/utils"

type CreateGenerationSettings struct {
	BasePrompt   utils.Nullable[string] `json:"basePrompt"`
	VarietyLevel utils.Nullable[int]    `json:"varietyLevel"`
}
