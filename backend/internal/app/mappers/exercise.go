package mappers

import (
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"
	desc "fitness-trainer/pkg/workouts"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ExerciseToProto(exercise domain.Exercise) *desc.Exercise {
	muscleGroups := make([]string, 0, len(exercise.TargetMuscleGroups))
	for _, muscleGroup := range exercise.TargetMuscleGroups {
		muscleGroups = append(muscleGroups, muscleGroup.String())
	}

	return &desc.Exercise{
		Id:                 exercise.ID.String(),
		Name:               exercise.Name,
		Description:        exercise.Description,
		TargetMuscleGroups: muscleGroups,
		CreatedAt:          timestamppb.New(exercise.CreatedAt),
		UpdatedAt:          timestamppb.New(exercise.UpdatedAt),
	}
}

func ExercisesToProto(exercises []domain.Exercise) []*desc.Exercise {
	result := make([]*desc.Exercise, 0, len(exercises))
	for _, exercise := range exercises {
		result = append(result, ExerciseToProto(exercise))
	}

	return result
}

func MuscleGroupDTOToProto(muscleGroup dto.MuscleGroupDTO) *desc.MuscleGroup {
	return &desc.MuscleGroup{
		Id:   muscleGroup.ID.String(),
		Name: muscleGroup.Name,
	}
}

func MuscleGroupDTOsToProto(muscleGroups []dto.MuscleGroupDTO) []*desc.MuscleGroup {
	result := make([]*desc.MuscleGroup, 0, len(muscleGroups))
	for _, muscleGroup := range muscleGroups {
		result = append(result, MuscleGroupDTOToProto(muscleGroup))
	}

	return result
}
