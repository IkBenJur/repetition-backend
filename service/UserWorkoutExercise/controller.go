package userWorkoutExercise

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

func (controller *Controller) SaveUserWorkoutExercise(workoutExercise types.UserWorkoutExercise) error {
	tx, err := controller.db.Begin()
	
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userWorkoutId int64
	err = tx.QueryRow("INSERT INTO userworkoutexercise (userworkoutid, exerciseid) VALUES ($1, $2) RETURNING id", workoutExercise.UserWorkoutId, workoutExercise.ExerciseId).Scan(&userWorkoutId)
	if err != nil {
		return err
	}

	exerciseSetStmt, err := tx.Prepare("INSERT INTO userworkoutexerciseset (userworkoutexerciseid, reps, weight) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return err
	}
	defer exerciseSetStmt.Close()

	for _, set := range workoutExercise.UserWorkoutExerciseSets {
		_, err := exerciseSetStmt.Exec(workoutExercise.ExerciseId, set.Reps, set.Weight)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (controller *Controller) FindUserIdForUserWorkoutExerciseId(id int) (int, error) {
	rows, err := controller.db.Query("SELECT uw.userid FROM userworkoutexercise uwe JOIN userworkout uw ON uw.id = uwe.userworkoutid WHERE uwe.id = $1 LIMIT 1", id)
	if err != nil {
		return 0, err
	}

	userId := 0
	for rows.Next() {
		err := rows.Scan(&userId)
		if err != nil {
			return 0, err
		}
	}

	return userId, nil
}