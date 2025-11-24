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

func (controller *Controller) SaveUserWorkoutExerciseSet(workoutExerciseSet types.UserWorkoutExerciseSet) error {
	_, err := controller.db.Exec("INSERT INTO userworkoutexerciseset (userworkoutexerciseid, reps, weight) VALUES ($1, $2, $3)", workoutExerciseSet.UserWorkoutExerciseId, workoutExerciseSet.Reps, workoutExerciseSet.Weight)
	
	return err
}