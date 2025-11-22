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
	for j, set := range payload.UserWorkoutExerciseSets {
		sets[j] = UserWorkoutExerciseSet{
			UserWorkoutExerciseId: set.UserWorkoutExerciseId,
			Reps:                  set.Reps,
			Weight:                set.Weight,
		}
	}

	return UserWorkoutExercise {
		UserWorkoutExerciseSets: sets,
		ExerciseId:              payload.ExerciseId,
	}
}