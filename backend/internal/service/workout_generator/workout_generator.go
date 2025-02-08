package workout_generator_service

import (
	"context"
	"encoding/json"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"
	"fmt"

	"github.com/opentracing/opentracing-go"
)

const (
	systemPrompt       = `Вы - всемирно известный фитнес-тренер. Ваша задача - проанализировать последние четыре тренировки вашего клиента и составить список упражнений для его новой тренировки. Убедитесь, что используются только упражнения из предоставленного списка. Необходимо, чтобы каждая мышечная группа, участвовавшая в последних четырех тренировках, была задействована. План тренировки должен включать упражнения как с использованием свободных весов, так и на тренажерах. Упражнения со свободным весом должны быть размещены в начале тренировки, так как они требуют больше энергии. Включите от 5 до 8 упражнений в тренировку. Если запрос пользователя включает дополнительные предпочтения, отдавайте приоритет их при выборе упражнений. При добавлении нескольких упражнений для одной мышечной группы в одну тренировку они должны быть разными упражнениями, а не вариациями одного и того же упражнения.`
	userPromptTemplate = `<exercise_list>%s</exercise_list>\n<workout_list>%s</workout_list>\n<user_preferences>%s</user_preferences>`
)

type CompletionProvider interface {
	CreateCompletion(ctx context.Context, userID domain.ID, systemPrompt, prompt string) (string, error)
}

type Service struct {
	completionProvider CompletionProvider
}

func New(completionProvider CompletionProvider) *Service {
	return &Service{
		completionProvider: completionProvider,
	}
}

func (s *Service) GenerateWorkout(ctx context.Context, userID domain.ID, workouts []dto.SlimWorkoutDTO, exercises []dto.SlimExerciseDTO, userPrompt string) (dto.GeneratedWorkoutDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "workout_generator_service.GenerateWorkout")
	defer span.Finish()

	marshaledWorkouts, err := marshalWorkouts(workouts)
	if err != nil {
		return dto.GeneratedWorkoutDTO{}, fmt.Errorf("failed to marshal workouts: %w", err)
	}

	marshaledExercises, err := marshalExercises(exercises)
	if err != nil {
		return dto.GeneratedWorkoutDTO{}, fmt.Errorf("failed to marshal exercises: %w", err)
	}

	innerUserPrompt := fmt.Sprintf(userPromptTemplate, marshaledExercises, marshaledWorkouts, userPrompt)

	completion, err := s.completionProvider.CreateCompletion(ctx, userID, systemPrompt, innerUserPrompt)
	if err != nil {
		return dto.GeneratedWorkoutDTO{}, fmt.Errorf("failed to create completion: %w", err)
	}

	return unmarshalCompletion(completion)
}

func marshalWorkouts(workouts []dto.SlimWorkoutDTO) (string, error) {
	type exercise struct {
		ID string `json:"id"`
	}

	type workout struct {
		ID        string     `json:"id"`
		CreatedAt string     `json:"created_at"`
		Exercises []exercise `json:"exercises"`
	}

	workoutsToMarshal := make([]workout, 0, len(workouts))
	for _, w := range workouts {
		exercises := make([]exercise, 0, len(w.ExerciseIDs))
		for _, e := range w.ExerciseIDs {
			exercises = append(exercises, exercise{ID: e.String()})
		}
		workoutsToMarshal = append(workoutsToMarshal, workout{
			ID:        w.ID.String(),
			CreatedAt: w.CreatedAt.String(),
			Exercises: exercises,
		})
	}

	return marshal(workoutsToMarshal)
}

func marshalExercises(exercises []dto.SlimExerciseDTO) (string, error) {
	type exercise struct {
		ID                 string   `json:"id"`
		Name               string   `json:"name"`
		TargetMuscleGroups []string `json:"targetMuscleGroups"`
	}

	exercisesToMarshal := make([]exercise, 0, len(exercises))
	for _, e := range exercises {
		exercisesToMarshal = append(exercisesToMarshal, exercise{
			ID:                 e.ID.String(),
			Name:               e.Name,
			TargetMuscleGroups: marshalMuscleGroups(e.TargetMuscleGroups),
		})
	}

	return marshal(exercisesToMarshal)
}

func marshalMuscleGroups(muscleGroups []domain.MuscleGroup) []string {
	muscleGroupsToMarshal := make([]string, 0, len(muscleGroups))
	for _, mg := range muscleGroups {
		muscleGroupsToMarshal = append(muscleGroupsToMarshal, mg.String())
	}

	return muscleGroupsToMarshal
}

func marshal(data interface{}) (string, error) {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(marshaledData), nil
}

func unmarshalCompletion(rawCompletion string) (dto.GeneratedWorkoutDTO, error) {
	type exercise struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	type completion struct {
		Exercises []exercise `json:"exercises"`
		Reasoning string     `json:"reasoning"`
	}

	var completionData completion
	err := json.Unmarshal([]byte(rawCompletion), &completionData)
	if err != nil {
		return dto.GeneratedWorkoutDTO{}, err
	}

	exerciseIDs := make([]domain.ID, 0, len(completionData.Exercises))
	for _, e := range completionData.Exercises {
		parsedID, err := domain.ParseID(e.ID)
		if err != nil {
			return dto.GeneratedWorkoutDTO{}, fmt.Errorf("failed to parse exercise ID: %w", err)
		}
		exerciseIDs = append(exerciseIDs, parsedID)
	}

	return dto.GeneratedWorkoutDTO{
		ExerciseIDs: exerciseIDs,
		Reasoning:   completionData.Reasoning,
	}, nil
}
