package userWorkout

import "github.com/IkBenJur/repetition-backend/types"

func UserWorkoutPayloadIntoUserWorkout(payload types.NewUserWorkoutPayload) types.UserWorkout {

	exercises := make([]types.UserWorkoutExercise, len(payload.UserWorkoutExercises))
	for i, exercise := range payload.UserWorkoutExercises {

		sets := make([]types.UserWorkoutExerciseSet, len(exercise.UserWorkoutExerciseSets))
		for j, set := range exercise.UserWorkoutExerciseSets {
			sets[j] = types.UserWorkoutExerciseSet{
				UserWorkoutExerciseId: set.UserWorkoutExerciseId,
				Reps:                  set.Reps,
				Weight:                set.Weight,
			}
		}

		exercises[i] = types.UserWorkoutExercise{
			UserWorkoutExerciseSets: sets,
			ExerciseId:              exercise.ExerciseId,
		}

	}

	return types.UserWorkout{
		Name:                 payload.Name,
		UserId:               payload.UserId,
		UserWorkoutExercises: exercises,
	}
}

