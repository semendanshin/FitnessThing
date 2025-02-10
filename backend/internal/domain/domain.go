package domain

import (
	"fitness-trainer/internal/utils"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ID uuid.UUID

func NewID() ID {
	return ID(uuid.New())
}

func (i ID) String() string {
	return uuid.UUID(i).String()
}

func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ID{}, err
	}
	return ID(id), nil
}

type Model struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewModel() Model {
	return Model{
		ID:        NewID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type User struct {
	Model

	Email         string
	Password      string
	FirstName     string
	LastName      string
	DateOfBirth   time.Time
	Height        float32
	Weight        float32
	ProfilePicURL string
}

func NewUser(
	Email string,
	Password string,
	FirstName string,
	LastName string,
	DateOfBirth time.Time,
	Height float32,
	Weight float32,
) User {
	return User{
		Model:       NewModel(),
		Email:       Email,
		Password:    Password,
		FirstName:   FirstName,
		LastName:    LastName,
		DateOfBirth: DateOfBirth,
		Height:      Height,
		Weight:      Weight,
	}
}

type MuscleGroup string

const (
	MuscleGroupChest      MuscleGroup = "chest"
	MuscleGroupBack       MuscleGroup = "lats"
	MuscleGroupShoulders  MuscleGroup = "shoulders"
	MuscleGroupBiceps     MuscleGroup = "biceps"
	MuscleGroupTriceps    MuscleGroup = "triceps"
	MuscleGroupForearms   MuscleGroup = "forearms"
	MuscleGroupAbs        MuscleGroup = "abs"
	MuscleGroupQuads      MuscleGroup = "quads"
	MuscleGroupHamstrings MuscleGroup = "hamstrings"
	MuscleGroupCalves     MuscleGroup = "calves"
	MuscleGroupGlutes     MuscleGroup = "glutes"
	MuscleGroupLowerBack  MuscleGroup = "lower-back"
	MuscleGroupTraps      MuscleGroup = "traps"
)

func (m MuscleGroup) String() string {
	return string(m)
}

func NewMuscleGroup(m string) (MuscleGroup, error) {
	switch m {
	case "chest":
		return MuscleGroupChest, nil
	case "lats":
		return MuscleGroupBack, nil
	case "shoulders":
		return MuscleGroupShoulders, nil
	case "biceps":
		return MuscleGroupBiceps, nil
	case "triceps":
		return MuscleGroupTriceps, nil
	case "forearms":
		return MuscleGroupForearms, nil
	case "abs":
		return MuscleGroupAbs, nil
	case "quads":
		return MuscleGroupQuads, nil
	case "hamstrings":
		return MuscleGroupHamstrings, nil
	case "calves":
		return MuscleGroupCalves, nil
	case "glutes":
		return MuscleGroupGlutes, nil
	case "lower back":
		return MuscleGroupLowerBack, nil
	case "traps":
		return MuscleGroupTraps, nil
	default:
		return "", fmt.Errorf("unknown muscle group: %w", ErrInvalidArgument)
	}
}

type Exercise struct {
	Model

	Name               string
	Description        string
	VideoURL           string
	TargetMuscleGroups []MuscleGroup
}

func NewExercise(name, description, videoURL string, targetMuscleGroups []MuscleGroup) Exercise {
	return Exercise{
		Model:              NewModel(),
		Name:               name,
		Description:        description,
		VideoURL:           videoURL,
		TargetMuscleGroups: targetMuscleGroups,
	}
}

type Routine struct {
	Model

	UserID      ID
	Name        string
	Description string
}

func NewRoutine(userID ID, name, description string) Routine {
	return Routine{
		Model:       NewModel(),
		UserID:      userID,
		Name:        name,
		Description: description,
	}
}

type SetType string

const (
	SetTypeUnknown SetType = ""
	SetTypeReps    SetType = "reps"
	SetTypeWeight  SetType = "weight"
	SetTypeTime    SetType = "time"
)

func (s SetType) String() string {
	return string(s)
}

func NewSetType(s string) (SetType, error) {
	switch s {
	case "reps":
		return SetTypeReps, nil
	case "weight":
		return SetTypeWeight, nil
	case "time":
		return SetTypeTime, nil
	default:
		return "", fmt.Errorf("unknown set type: %w", ErrInvalidArgument)
	}
}

type ExerciseInstance struct {
	Model

	RoutineID  ID
	ExerciseID ID
}

func NewExerciseInstance(routineID, exerciseID ID) ExerciseInstance {
	return ExerciseInstance{
		Model:      NewModel(),
		RoutineID:  routineID,
		ExerciseID: exerciseID,
	}
}

type Set struct {
	Model

	ExerciseInstanceID ID
	SetType            SetType
	Reps               int
	Weight             float32
	Time               time.Duration
}

func NewSet(exerciseInstanceID ID, setType SetType, reps int, weight float32, time time.Duration) Set {
	return Set{
		Model:              NewModel(),
		ExerciseInstanceID: exerciseInstanceID,
		SetType:            setType,
		Reps:               reps,
		Weight:             weight,
		Time:               time,
	}
}

type Workout struct {
	Model

	UserID        ID
	RoutineID     utils.Nullable[ID]
	Notes         string
	Rating        int
	FinishedAt    time.Time
	IsAIGenerated bool
	Reasoning     string
}

func NewWorkout(userID ID, routineID utils.Nullable[ID], isAIGenerated bool) Workout {
	return Workout{
		Model:         NewModel(),
		UserID:        userID,
		RoutineID:     routineID,
		IsAIGenerated: isAIGenerated,
	}
}

type ExerciseLog struct {
	Model

	WorkoutID   ID
	ExerciseID  ID
	Notes       string
	PowerRating int
}

func NewExerciseLog(workoutID, exerciseID ID) ExerciseLog {
	return ExerciseLog{
		Model:      NewModel(),
		WorkoutID:  workoutID,
		ExerciseID: exerciseID,
	}
}

type ExpectedSet struct {
	Model

	ExerciseLogID ID
	SetType       SetType
	Reps          int
	Weight        float32
	Time          time.Duration
}

func NewExpectedSet(exerciseLogID ID, setType SetType, reps int, weight float32, time time.Duration) ExpectedSet {
	return ExpectedSet{
		Model:         NewModel(),
		ExerciseLogID: exerciseLogID,
		SetType:       setType,
		Reps:          reps,
		Weight:        weight,
		Time:          time,
	}
}

type ExerciseSetLog struct {
	Model

	ExerciseLogID ID
	Reps          int
	Weight        float32
	Time          time.Duration
}

func NewExerciseSetLog(exerciseLogID ID, reps int, weight float32, time time.Duration) ExerciseSetLog {
	return ExerciseSetLog{
		Model:         NewModel(),
		ExerciseLogID: exerciseLogID,
		Reps:          reps,
		Weight:        weight,
		Time:          time,
	}
}

type Session struct {
	Model

	UserID    ID
	ExpiredAt time.Time
	Token     string
}

func NewSession(userID ID, expiredAt time.Time, token string) Session {
	return Session{
		Model:     NewModel(),
		UserID:    userID,
		ExpiredAt: expiredAt,
		Token:     token,
	}
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

/*
- Получить список планов тренировок пользователя и выбрать план тренировки или Автоматически создать план для сегодняшней тренировки
	- Просмотреть план тренировки
		- Просмотреть информацию об упражнении (название, описание, техника выполнения, видео)
		- Просмотреть список подходов к упражнению
	- Редактировать план тренировки
		- Добавить упражнение в план тренировки
			- Просмотреть список упражнений
		- Удалить упражнение из плана тренировки
		- Редактировать упражнение в плане тренировки
			- Добавить подход к упражнению в плане тренировки
			- Удалить подход к упражнению из плана тренировки
			- Редактировать подход к упражнению в плане тренировки
		- Заменить упражнение в плане тренировки на аналогичное
			- Просмотреть список аналогичных упражнений

	- Начать тренировку
		- Создать тренировку
			- Просмотреть список планов тренировок пользователя
		- Добавить отчет о выполнении упражнения
			- Добавить отчет о выполнении подхода к упражнению
			- Оценить силовую нагрузку
			- Добавить заметки
		- Завершить тренировку
			- Просмотреть отчет о выполнении тренировки
			- Оценить тренировку
*/
