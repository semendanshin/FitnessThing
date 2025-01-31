package mappers

import (
	"fitness-trainer/internal/domain"
	desc "fitness-trainer/pkg/workouts"
)

func SetTypeToProto(setType domain.SetType) desc.SetType {
	switch setType {
	case domain.SetTypeReps:
		return desc.SetType_SET_TYPE_REPS
	case domain.SetTypeWeight:
		return desc.SetType_SET_TYPE_WEIGHT
	case domain.SetTypeTime:
		return desc.SetType_SET_TYPE_TIME
	default:
		return desc.SetType_SET_TYPE_UNSPECIFIED
	}
}

func SetTypeFromProto(setType desc.SetType) domain.SetType {
	switch setType {
	case desc.SetType_SET_TYPE_REPS:
		return domain.SetTypeReps
	case desc.SetType_SET_TYPE_WEIGHT:
		return domain.SetTypeWeight
	case desc.SetType_SET_TYPE_TIME:
		return domain.SetTypeTime
	default:
		return domain.SetTypeUnknown
	}
}
