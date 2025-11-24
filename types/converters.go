package types

func UserWorkoutPayloadIntoUserWorkout(payload NewUserWorkoutPayload) UserWorkout {
	exercises := make([]UserWorkoutExercise, len(payload.UserWorkoutExercises))
	for i, exercise := range payload.UserWorkoutExercises {
		exercises[i] = UserWorkoutexerciseIntoUserWorkoutexercise(exercise)
	}

	return UserWorkout{
		Name:                 payload.Name,
		UserId:               payload.UserId,
		UserWorkoutExercises: exercises,
	}
}

func UserWorkoutexerciseIntoUserWorkoutexercise(payload UserWorkoutExercisePayload) UserWorkoutExercise {
	sets := make([]UserWorkoutExerciseSet, len(payload.UserWorkoutExerciseSets))
	for i, set := range payload.UserWorkoutExerciseSets {
		sets[i] = UserWorkoutExerciseSetIntoUserWorkoutExerciseSet(set)
	}

	return UserWorkoutExercise {
		UserWorkoutExerciseSets: sets,
		ExerciseId:              payload.ExerciseId,
		UserWorkoutId:           payload.UserWorkoutId,
	}
}

func UserWorkoutExerciseSetIntoUserWorkoutExerciseSet(payload UserWorkoutExerciseSetPayload) UserWorkoutExerciseSet {
	return UserWorkoutExerciseSet {
			UserWorkoutExerciseId: payload.UserWorkoutExerciseId,
			Reps:                  payload.Reps,
			Weight:                payload.Weight,
	}
}