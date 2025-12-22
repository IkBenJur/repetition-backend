package types

func UserWorkoutPayloadIntoUserWorkout(payload NewUserWorkoutPayload) UserWorkout {
	exercises := make([]*UserWorkoutExercise, len(payload.UserWorkoutExercises))
	for i, exercise := range payload.UserWorkoutExercises {
		exercises[i] = UserWorkoutExercisePayloadIntoUserWorkoutExercise(exercise)
	}

	return UserWorkout{
		Name:                 payload.Name,
		UserId:               payload.UserId,
		UserWorkoutExercises: exercises,
	}
}

func UserWorkoutExercisePayloadIntoUserWorkoutExercise(payload UserWorkoutExercisePayload) *UserWorkoutExercise {
	sets := make([]*UserWorkoutExerciseSet, len(payload.UserWorkoutExerciseSets))
	for i, set := range payload.UserWorkoutExerciseSets {
		sets[i] = set.ToEntity()
	}

	return &UserWorkoutExercise{
		UserWorkoutExerciseSets: sets,
		ExerciseId:              payload.ExerciseId,
		ExerciseName:            payload.ExerciseName,
		UserWorkoutId:           payload.UserWorkoutId,
	}
}
