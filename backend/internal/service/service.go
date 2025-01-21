package service

import (
	"context"
	"time"

	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"
)

type jwtProvider interface {
	GeneratePair(ctx context.Context, userID, pairID domain.ID, atTime time.Time) (domain.Tokens, error)
	VerifyPair(ctx context.Context, userID domain.ID, tokens domain.Tokens, atTime time.Time) error
	ParseToken(ctx context.Context, token string) (domain.ID, error)
}

type sessionRepository interface {
	GetSessionByToken(ctx context.Context, token string) (domain.Session, error)
	SetSessionExpired(ctx context.Context, id domain.ID, expiredAt time.Time) error
	CreateSession(ctx context.Context, session domain.Session) (domain.Session, error)
}

type userRepository interface {
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id domain.ID) (domain.User, error)
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
}

type exerciseRepository interface {
	GetExercises(ctx context.Context, muscleGroups, excludedExercises []domain.ID) ([]domain.Exercise, error)
	GetExerciseByID(ctx context.Context, id domain.ID) (domain.Exercise, error)
	CreateExercise(ctx context.Context, exercise domain.Exercise) (domain.Exercise, error)
}

type routineRepository interface {
	GetRoutines(ctx context.Context, userID domain.ID) ([]domain.Routine, error)
	CreateRoutine(ctx context.Context, routine domain.Routine) (domain.Routine, error)
	GetRoutineByID(ctx context.Context, id domain.ID) (domain.Routine, error)
	DeleteRoutine(ctx context.Context, id domain.ID) error
}

type exerciseInstanceRepository interface {
	GetExerciseInstancesByRoutineID(ctx context.Context, routineID domain.ID) ([]domain.ExerciseInstance, error)
	CreateExerciseInstance(ctx context.Context, exerciseInstance domain.ExerciseInstance) (domain.ExerciseInstance, error)
	DeleteExerciseInstance(ctx context.Context, id domain.ID) error
}

type muscleGroupRepository interface {
	GetMuscleGroups(ctx context.Context) ([]dto.MuscleGroupDTO, error)
	GetMuscleGroupByName(ctx context.Context, name string) (dto.MuscleGroupDTO, error)
}

type workoutRepository interface {
	CreateWorkout(ctx context.Context, workout domain.Workout) (domain.Workout, error)
	GetWorkoutByID(ctx context.Context, id domain.ID) (domain.Workout, error)
	GetActiveWorkouts(ctx context.Context, userID domain.ID) ([]domain.Workout, error)
	UpdateWorkout(ctx context.Context, id domain.ID, workout domain.Workout) (domain.Workout, error)
}

type exerciseLogRepository interface {
	GetExerciseLogsByWorkoutID(ctx context.Context, workoutID domain.ID) ([]domain.ExerciseLog, error)
	GetExerciseLogByID(ctx context.Context, id domain.ID) (domain.ExerciseLog, error)
	CreateExerciseLog(ctx context.Context, exerciseLog domain.ExerciseLog) (domain.ExerciseLog, error)
	GetExerciseLogsByExerciseIDAndUserID(ctx context.Context, exerciseID, userID domain.ID) ([]domain.ExerciseLog, error)
}

type setLogRepository interface {
	GetSetLogsByExerciseLogID(ctx context.Context, exerciseLogID domain.ID) ([]domain.ExerciseSetLog, error)
	CreateSetLog(ctx context.Context, setLog domain.ExerciseSetLog) (domain.ExerciseSetLog, error)
}

type Service struct {
	jwtProvider                jwtProvider
	sessionRepository          sessionRepository
	userRepository             userRepository
	exerciseRepository         exerciseRepository
	routineRepository          routineRepository
	exerciseInstanceRepository exerciseInstanceRepository
	muscleGroupRepository      muscleGroupRepository
	workoutRepository          workoutRepository
	exerciseLogRepository      exerciseLogRepository
	setLogRepository           setLogRepository
}

func New(
	jwtProvider jwtProvider,
	sessionRepository sessionRepository,
	userRepository userRepository,
	exerciseRepository exerciseRepository,
	routineRepository routineRepository,
	exerciseInstanceRepository exerciseInstanceRepository,
	muscleGroupRepository muscleGroupRepository,
	workoutRepository workoutRepository,
	exerciseLogRepository exerciseLogRepository,
	setLogRepository setLogRepository,
) *Service {
	return &Service{
		jwtProvider:                jwtProvider,
		sessionRepository:          sessionRepository,
		userRepository:             userRepository,
		exerciseRepository:         exerciseRepository,
		routineRepository:          routineRepository,
		exerciseInstanceRepository: exerciseInstanceRepository,
		muscleGroupRepository:      muscleGroupRepository,
		workoutRepository:          workoutRepository,
		exerciseLogRepository:      exerciseLogRepository,
		setLogRepository:           setLogRepository,
	}
}
