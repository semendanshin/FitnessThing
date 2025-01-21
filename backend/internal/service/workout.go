package service

import (
	"context"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"
	"fitness-trainer/internal/logger"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
)

func (s *Service) StartWorkout(ctx context.Context, userID domain.ID, routineID *domain.ID) (domain.Workout, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.StartWorkout")
	defer span.Finish()

	if routineID != nil {
		_, err := s.routineRepository.GetRoutineByID(ctx, *routineID)
		if err != nil {
			return domain.Workout{}, fmt.Errorf("%w: %w", domain.ErrInvalidArgument, err)
		}
	}

	workout := domain.NewWorkout(userID, routineID)

	workout, err := s.workoutRepository.CreateWorkout(ctx, workout)
	if err != nil {
		return domain.Workout{}, err
	}

	if routineID != nil {
		err = s.assignExercisesToWorkout(ctx, workout)
		if err != nil {
			return domain.Workout{}, err
		}
	}

	return workout, nil
}

func (s *Service) assignExercisesToWorkout(ctx context.Context, workout domain.Workout) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.assignExercisesToWorkout")
	defer span.Finish()

	routine, err := s.routineRepository.GetRoutineByID(ctx, workout.RoutineID.V)
	if err != nil {
		return err
	}

	exerciseInstances, err := s.exerciseInstanceRepository.GetExerciseInstancesByRoutineID(ctx, routine.ID)
	if err != nil {
		return err
	}

	for _, instance := range exerciseInstances {
		s.LogExercise(ctx, workout.UserID, workout.ID, instance.ExerciseID)
	}

	return nil
}

func (s *Service) GetWorkout(ctx context.Context, userID domain.ID, workoutID domain.ID) (dto.WorkoutDetailsDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetWorkout")
	defer span.Finish()

	workout, err := s.workoutRepository.GetWorkoutByID(ctx, workoutID)
	if err != nil {
		return dto.WorkoutDetailsDTO{}, err
	}

	if workout.UserID != userID {
		logger.Errorf("user %s tried to access workout %s", userID, workoutID)
		return dto.WorkoutDetailsDTO{}, domain.ErrNotFound
	}

	exerciseLogs, err := s.exerciseLogRepository.GetExerciseLogsByWorkoutID(ctx, workoutID)
	if err != nil {
		return dto.WorkoutDetailsDTO{}, err
	}

	exerciseLogsDTOs := make([]dto.ExerciseLogDTO, 0, len(exerciseLogs))
	for _, exerciseLog := range exerciseLogs {
		exerciseLogDTO, err := s.GetExerciseLog(ctx, userID, exerciseLog.ID)
		if err != nil {
			return dto.WorkoutDetailsDTO{}, err
		}

		exerciseLogsDTOs = append(exerciseLogsDTOs, exerciseLogDTO)
	}

	return dto.WorkoutDetailsDTO{
		Workout:      workout,
		ExerciseLogs: exerciseLogsDTOs,
	}, nil
}

func (s *Service) GetExerciseLog(ctx context.Context, userID, exerciseLogID domain.ID) (dto.ExerciseLogDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetExerciseLog")
	defer span.Finish()

	exerciseLog, err := s.exerciseLogRepository.GetExerciseLogByID(ctx, exerciseLogID)
	if err != nil {
		return dto.ExerciseLogDTO{}, err
	}

	workout, err := s.workoutRepository.GetWorkoutByID(ctx, exerciseLog.WorkoutID)
	if err != nil {
		return dto.ExerciseLogDTO{}, err
	}

	if workout.UserID != userID {
		logger.Errorf("user %s tried to access exercise log %s", userID, exerciseLogID)
		return dto.ExerciseLogDTO{}, domain.ErrNotFound
	}

	setLogs, err := s.setLogRepository.GetSetLogsByExerciseLogID(ctx, exerciseLogID)
	if err != nil {
		return dto.ExerciseLogDTO{}, err
	}

	exercise, err := s.exerciseRepository.GetExerciseByID(ctx, exerciseLog.ExerciseID)
	if err != nil {
		return dto.ExerciseLogDTO{}, err
	}

	return dto.ExerciseLogDTO{
		ExerviceLog: exerciseLog,
		SetLogs:     setLogs,
		Exercise:    exercise,
	}, nil
}

func (s *Service) LogExercise(ctx context.Context, userID, workoutID, exerciseID domain.ID) (domain.ExerciseLog, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.LogExercise")
	defer span.Finish()

	workout, err := s.workoutRepository.GetWorkoutByID(ctx, workoutID)
	if err != nil {
		return domain.ExerciseLog{}, err
	}

	if !workout.FinishedAt.IsZero() {
		logger.Errorf("user %s tried to log exercise for finished workout %s", userID, workoutID)
		return domain.ExerciseLog{}, fmt.Errorf("%w: workout %s is already finished", domain.ErrInvalidArgument, workoutID)
	}

	_, err = s.exerciseRepository.GetExerciseByID(ctx, exerciseID)
	if err != nil {
		return domain.ExerciseLog{}, err
	}

	if workout.UserID != userID {
		logger.Errorf("user %s tried to log exercise for workout %s", userID, workoutID)
		return domain.ExerciseLog{}, domain.ErrNotFound
	}

	exerciseLog := domain.NewExerciseLog(workoutID, exerciseID)

	exerciseLog, err = s.exerciseLogRepository.CreateExerciseLog(ctx, exerciseLog)
	if err != nil {
		return domain.ExerciseLog{}, err
	}

	return exerciseLog, nil
}

func (s *Service) LogSet(ctx context.Context, userID, workoutID, exerciseLogID domain.ID, setlogDTO dto.CreateSetLogDTO) (domain.ExerciseSetLog, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.LogSet")
	defer span.Finish()

	exerciseLog, err := s.exerciseLogRepository.GetExerciseLogByID(ctx, exerciseLogID)
	if err != nil {
		return domain.ExerciseSetLog{}, err
	}

	workout, err := s.workoutRepository.GetWorkoutByID(ctx, workoutID)
	if err != nil {
		return domain.ExerciseSetLog{}, err
	}

	if !workout.FinishedAt.IsZero() {
		logger.Errorf("user %s tried to log exercise for finished workout %s", userID, workoutID)
		return domain.ExerciseSetLog{}, fmt.Errorf("%w: workout %s is already finished", domain.ErrInvalidArgument, workoutID)
	}

	if workout.UserID != userID {
		logger.Errorf("user %s tried to log set for workout %s", userID, workoutID)
		return domain.ExerciseSetLog{}, domain.ErrNotFound
	}

	if exerciseLog.WorkoutID != workoutID {
		logger.Errorf("user %s tried to log set for exercise log %s for workout %s", userID, exerciseLogID, workoutID)
		return domain.ExerciseSetLog{}, domain.ErrNotFound
	}

	setLog := domain.NewExerciseSetLog(
		exerciseLogID,
		setlogDTO.Reps,
		setlogDTO.Weight,
		time.Duration(0),
	)

	setLog, err = s.setLogRepository.CreateSetLog(ctx, setLog)
	if err != nil {
		return domain.ExerciseSetLog{}, err
	}

	return setLog, nil
}

func (s *Service) GetActiveWorkouts(ctx context.Context, userID domain.ID) ([]domain.Workout, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetActiveWorkouts")
	defer span.Finish()

	workouts, err := s.workoutRepository.GetActiveWorkouts(ctx, userID)
	if err != nil {
		return nil, err
	}

	return workouts, nil
}

func (s *Service) CompleteWorkout(ctx context.Context, userID, workoutID domain.ID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.FinishWorkout")
	defer span.Finish()

	workout, err := s.workoutRepository.GetWorkoutByID(ctx, workoutID)
	if err != nil {
		return err
	}

	if workout.UserID != userID {
		logger.Errorf("user %s tried to finish workout %s", userID, workoutID)
		return domain.ErrNotFound
	}

	if !workout.FinishedAt.IsZero() {
		logger.Errorf("user %s tried to finish already finished workout %s", userID, workoutID)
		return fmt.Errorf("%w: workout %s is already finished", domain.ErrInvalidArgument, workoutID)
	}

	workout.FinishedAt = time.Now()

	_, err = s.workoutRepository.UpdateWorkout(ctx, workoutID, workout)
	if err != nil {
		return err
	}

	return nil
}
