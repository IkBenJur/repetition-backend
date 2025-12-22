package userWorkoutExerciseSet

import (
	"database/sql"

	"github.com/IkBenJur/repetition-backend/types"
)

type Controller struct {
	db *sql.DB
}

func NewController(db *sql.DB) *Controller {
	return &Controller{db: db}
}

func (controller *Controller) CreateNewUserWorkoutExerciseSet(workoutExerciseSet types.UserWorkoutExerciseSet) (int, error) {
	var workoutExerciseSetId int
	err := controller.db.QueryRow("INSERT INTO userworkoutexerciseset (userworkoutexerciseid, reps, weight) VALUES ($1, $2, $3) RETURNING id", workoutExerciseSet.UserWorkoutExerciseId, workoutExerciseSet.Reps, workoutExerciseSet.Weight).Scan(&workoutExerciseSetId)

	return workoutExerciseSetId, err
}
	
func (controller *Controller) FindUserIdForSetId(workoutExerciseSet types.UserWorkoutExerciseSet) (int, error) {
	var userId int
	err := controller.db.QueryRow(`
		SELECT workout.userid FROM userworkoutexerciseset exerciseset
			JOIN userworkoutexercise exercise ON exercise.id = exerciseset.userworkoutexerciseid
			JOIN userworkout workout ON workout.id = exercise.userworkoutid
		WHERE exerciseset.id = $1`, workoutExerciseSet.ID).Scan(&userId)

	return userId, err
}

func (controller *Controller) UpdateUserWorkoutExerciseSet(workoutExerciseSet types.UserWorkoutExerciseSet) error {
	_, err := controller.db.Exec("UPDATE userworkoutexerciseset SET reps = $1, weight = $2 WHERE id = $3", workoutExerciseSet.Reps, workoutExerciseSet.Weight, workoutExerciseSet.ID)
	return err
}
