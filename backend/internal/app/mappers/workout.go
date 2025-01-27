package mappers

import (
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/domain/dto"
	desc "fitness-trainer/pkg/workouts"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func WorkoutsToProto(workouts []domain.Workout) *desc.WorkoutsListResponse {
	workoutsList := make([]*desc.Workout, 0, len(workouts))
	for _, workout := range workouts {
		workoutsList = append(workoutsList, WorkoutToProto(workout))
	}

	return &desc.WorkoutsListResponse{
		Workouts: workoutsList,
	}
}

func WorkoutToProto(workout domain.Workout) *desc.Workout {
	var routineID *string
	if workout.RoutineID.IsValid {
		routineIDValue := workout.RoutineID.V.String()
		routineID = &routineIDValue
	}
	return &desc.Workout{
		Id:         workout.ID.String(),
		RoutineId:  routineID,
		UserId:     workout.UserID.String(),
		CreatedAt:  timestamppb.New(workout.CreatedAt),
		Notes:      workout.Notes,
		Rating:     int32(workout.Rating),
		FinishedAt: timestamppb.New(workout.FinishedAt),
		UpdatedAt:  timestamppb.New(workout.UpdatedAt),
	}
}

func SetLogToProto(setLog domain.ExerciseSetLog) *desc.SetLog {
	return &desc.SetLog{
		Id:        setLog.ID.String(),
		Reps:      int32(setLog.Reps),
		Weight:    setLog.Weight,
		CreatedAt: timestamppb.New(setLog.CreatedAt),
		UpdatedAt: timestamppb.New(setLog.UpdatedAt),
	}
}

func ExerciseLogToProto(exerciseLog domain.ExerciseLog) *desc.ExerciseLog {
	return &desc.ExerciseLog{
		Id:         exerciseLog.ID.String(),
		WorkoutId:  exerciseLog.WorkoutID.String(),
		ExerciseId: exerciseLog.ExerciseID.String(),
		CreatedAt:  timestamppb.New(exerciseLog.CreatedAt),
		UpdatedAt:  timestamppb.New(exerciseLog.UpdatedAt),
	}
}

func ExerciseLogDTOsToProto(in []dto.ExerciseLogDTO) []*desc.ExerciseLogDetails {
	out := make([]*desc.ExerciseLogDetails, 0, len(in))
	for _, ex := range in {
		out = append(out, ExerciseLogDTOToProto(ex))
	}
	return out
}

func ExerciseLogDTOToProto(in dto.ExerciseLogDTO) *desc.ExerciseLogDetails {
	return &desc.ExerciseLogDetails{
		ExerciseLog: ExerciseLogToProto(in.ExerviceLog),
		Exercise:    ExerciseToProto(in.Exercise),
		SetLogs:     SetLogsToProto(in.SetLogs),
	}
}

func SetLogsToProto(setLogs []domain.ExerciseSetLog) []*desc.SetLog {
	setLogsList := make([]*desc.SetLog, 0, len(setLogs))
	for _, setLog := range setLogs {
		setLogsList = append(setLogsList, SetLogToProto(setLog))
	}

	return setLogsList
}

func ExerciseLogsToProto(exerciseLogs []domain.ExerciseLog) []*desc.ExerciseLog {
	result := make([]*desc.ExerciseLog, 0, len(exerciseLogs))
	for _, exerciseLog := range exerciseLogs {
		result = append(result, ExerciseLogToProto(exerciseLog))
	}

	return result
}

func WorkoutDTOToProto(workoutDTO dto.WorkoutDTO) *desc.GetWorkoutsResponse_WorkoutDetails {
	return &desc.GetWorkoutsResponse_WorkoutDetails{
		Workout:      WorkoutToProto(workoutDTO.Workout),
		ExerciseLogs: ExerciseLogsToProto(workoutDTO.ExerciseLogs),
	}
}

func WorkoutsDTOToProto(workoutDTOs []dto.WorkoutDTO) []*desc.GetWorkoutsResponse_WorkoutDetails {
	result := make([]*desc.GetWorkoutsResponse_WorkoutDetails, 0, len(workoutDTOs))
	for _, workoutDTO := range workoutDTOs {
		result = append(result, WorkoutDTOToProto(workoutDTO))
	}

	return result
}
