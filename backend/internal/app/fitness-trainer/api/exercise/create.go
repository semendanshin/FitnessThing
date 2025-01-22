package exercise

// import (
// 	"context"
// 	"fitness-trainer/internal/app/mappers"
// 	"fitness-trainer/internal/domain"
// 	"fitness-trainer/internal/logger"
// 	desc "fitness-trainer/pkg/workouts"

// 	"github.com/opentracing/opentracing-go"
// )

// func (i *Implementation) CreateExercise(ctx context.Context, in *desc.CreateExerciseRequest) (*desc.ExerciseResponse, error) {
// 	span, ctx := opentracing.StartSpanFromContext(ctx, "api.CreateExercise")
// 	defer span.Finish()

// 	var exerciseDTO dto.Exe

// 	exercise, err := i.service.CreateExercise(
// 		ctx,
		
// 	)
// 	if err != nil {
// 		logger.Errorf("error creating exercise: %v", err)
// 		return nil, err
// 	}

// 	return &desc.CreateExerciseResponse{
// 		Exercise: mappers.ExerciseToProto(exercise),
// 	}, nil
// }