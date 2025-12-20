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
