package mappers

import (
	"fitness-trainer/internal/domain"
	desc "fitness-trainer/pkg/workouts"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func SetToProto(set domain.Set) *desc.Set {
	return &desc.Set{
		Id:                 set.ID.String(),
		ExerciseInstanceId: set.ExerciseInstanceID.String(),
		SetType:            SetTypeToProto(set.SetType),
		Reps:               int32(set.Reps),
		Weight:             set.Weight,
		Time:               durationpb.New(set.Time),
		CreatedAt:          timestamppb.New(set.CreatedAt),
		UpdatedAt:          timestamppb.New(set.UpdatedAt),
	}
}

func SetsToProto(sets []domain.Set) []*desc.Set {
	result := make([]*desc.Set, 0, len(sets))
	for _, set := range sets {
		result = append(result, SetToProto(set))
	}

	return result
}
