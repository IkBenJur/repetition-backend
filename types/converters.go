package types

func UserWorkoutPayloadIntoUserWorkout(payload NewUserWorkoutPayload) UserWorkout {
	exercises := make([]*UserWorkoutExercise, len(payload.UserWorkoutExercises))
	for i, exercise := range payload.UserWorkoutExercises {
		exercises[i] = exercise.ToEntity()
	}

	return UserWorkout{
		Name:                 payload.Name,
		UserId:               payload.UserId,
		UserWorkoutExercises: exercises,
	}
}
