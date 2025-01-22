package service

import (
	"context"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"

	"github.com/opentracing/opentracing-go"
)

func (s *Service) GetRoutines(ctx context.Context, userID domain.ID) ([]domain.Routine, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetRoutines")
	defer span.Finish()

	return s.routineRepository.GetRoutines(ctx, userID)
}

func (s *Service) CreateRoutine(ctx context.Context, dto dto.CreateRoutineDTO) (domain.Routine, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.CreateRoutine")
	defer span.Finish()

	routine := domain.NewRoutine(dto.UserID, dto.Name, dto.Description)

	routine, err := s.routineRepository.CreateRoutine(ctx, routine)
	if err != nil {
		return domain.Routine{}, err
	}

	if dto.WorkoutID.IsValid {
		s.fillRoutineWithWorkout(ctx, routine, dto.WorkoutID.V)
	}

	return routine, nil
}

func (s *Service) fillRoutineWithWorkout(ctx context.Context, routine domain.Routine, workoutID domain.ID) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.fillRoutineWithWorkout")
	defer span.Finish()

	workout, err := s.workoutRepository.GetWorkoutByID(ctx, workoutID)
	if err != nil {
		return
	}

	if routine.UserID != workout.UserID {
		return
	}

	exerciseLogs, err := s.exerciseLogRepository.GetExerciseLogsByWorkoutID(ctx, workoutID)
	if err != nil {
		return
	}

	for _, exerciseLog := range exerciseLogs {
		s.AddExerciseToRoutine(ctx, routine.ID, exerciseLog.ExerciseID)

		// TODO: add sets and reps
	}
}

func (s *Service) GetRoutineByID(ctx context.Context, id domain.ID) (dto.RoutineDetailsDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetRoutineByID")
	defer span.Finish()

	routine, err := s.routineRepository.GetRoutineByID(ctx, id)
	if err != nil {
		return dto.RoutineDetailsDTO{}, err
	}

	exerciseInstances, err := s.exerciseInstanceRepository.GetExerciseInstancesByRoutineID(ctx, id)
	if err != nil {
		return dto.RoutineDetailsDTO{}, err
	}

	result := dto.RoutineDetailsDTO{
		ID:                routine.ID,
		UserID:            routine.UserID,
		Name:              routine.Name,
		Description:       routine.Description,
		CreatedAt:         routine.CreatedAt,
		UpdatedAt:         routine.UpdatedAt,
		ExerciseInstances: make([]dto.ExerciseInstanceDetailsDTO, len(exerciseInstances)),
	}

	for i, instance := range exerciseInstances {
		exercise, err := s.exerciseRepository.GetExerciseByID(ctx, instance.ExerciseID)
		if err != nil {
			return dto.RoutineDetailsDTO{}, err
		}

		// TODO: add sets

		result.ExerciseInstances[i] = dto.ExerciseInstanceDetailsDTO{
			ID:         instance.ID,
			RoutineID:  instance.RoutineID,
			ExerciseID: instance.ExerciseID,
			CreatedAt:  instance.CreatedAt,
			UpdatedAt:  instance.UpdatedAt,
			Exercise:   exercise,
			Sets:       []domain.Set{},
		}
	}

	return result, nil
}

func (s *Service) AddExerciseToRoutine(ctx context.Context, routineID domain.ID, exerciseID domain.ID) (domain.ExerciseInstance, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.AddExerciseToRoutine")
	defer span.Finish()

	exerciseInstance := domain.NewExerciseInstance(routineID, exerciseID)
	return s.exerciseInstanceRepository.CreateExerciseInstance(ctx, exerciseInstance)
}

func (s *Service) RemoveExerciseInstanceFromRoutine(ctx context.Context, userID domain.ID, routineID domain.ID, exerciseInstanceID domain.ID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.RemoveExerciseInstanceFromRoutine")
	defer span.Finish()

	routine, err := s.routineRepository.GetRoutineByID(ctx, routineID)
	if err != nil {
		return err
	}

	if routine.UserID != userID {
		return domain.ErrUnauthorized
	}

	return s.exerciseInstanceRepository.DeleteExerciseInstance(ctx, exerciseInstanceID)
}

func (s *Service) DeleteRoutine(ctx context.Context, id domain.ID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.DeleteRoutine")
	defer span.Finish()

	return s.routineRepository.DeleteRoutine(ctx, id)
}

func (s *Service) UpdateRoutine(ctx context.Context, id domain.ID, dto dto.UpdateRoutineDTO) (domain.Routine, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.UpdateRoutine")
	defer span.Finish()

	routine, err := s.routineRepository.GetRoutineByID(ctx, id)
	if err != nil {
		return domain.Routine{}, err
	}

	if dto.Name.IsValid {
		routine.Name = dto.Name.V
	}

	if dto.Description.IsValid {
		routine.Description = dto.Description.V
	}

	return s.routineRepository.UpdateRoutine(ctx, id, routine)
}
