package userWorkout

import (
	"database/sql"

	"github.com/IkBenJur/repetition-backend/types"
)

type Controller struct {
	db *sql.DB
}

func NewController (db *sql.DB) *Controller {
	return &Controller{ db: db }
}

func (controller *Controller) SaveUserWorkout(workout types.UserWorkout) error {
	tx, err := controller.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var workoutId int64
	err = tx.QueryRow("INSERT INTO userWorkout (name) VALUES ($1)", workout.Name).Scan(&workoutId)
	if err != nil {
		return err
	}
	
	exerciseStmt, err := tx.Prepare("INSERT INTO userworkoutexercise (userworkoutid, exerciseid) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer exerciseStmt.Close()

	exerciseSetStmt, err := tx.Prepare("INSERT INTO userworkoutexerciseset (userworkoutexerciseid, reps, weight) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer exerciseSetStmt.Close()

	for _, exercise := range workout.UserWorkoutExercises {
		var exerciseId int64
		err = exerciseStmt.QueryRow(workoutId, exercise.ExerciseId).Scan(&exerciseId)
		if err != nil {
			return err
		}
		
		for _, set := range exercise.UserWorkoutExerciseSets {
			_, err := exerciseSetStmt.Exec(exerciseId, set.Reps, set.Weight)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}