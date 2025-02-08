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
	systemPrompt = `
		Ты профессиональный всемирно известный фитнес-тренер. Тебе предстоит проанализировать 4 последние тренировки твоего клиента и составить список упражнений для новой тренировки. Обязательно используй только упражнения из представленного списка. При составлении программы нужно обеспечить, чтобы за последние 4 тренировки была проработана каждая группа мышц в теле, причём тренировочный план должен включать упражнения как со свободными весами, так и на тренажёрах. Упражнения со свободными весами ставь в начало тренировки, так как они более энергозатратны. В тренировку нужно включить не менее 5 и не более 8 упражнений. Если в пользовательском запросе присутствуют дополнительные пожелания, обязательно учти их при выборе упражнений, отдавая им более высокий приоритет. Лучше добавить разнообразия в тренировки, сохраняя их эффективность.
	`
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
