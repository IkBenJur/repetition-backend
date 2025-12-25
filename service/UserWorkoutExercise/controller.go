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

func (controller *Controller) CreateNewUserWorkoutExercise(workoutExercise types.UserWorkoutExercise) (int, error) {
	tx, err := controller.db.Begin()

	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var userWorkoutExerciseId int
	err = tx.QueryRow("INSERT INTO userworkoutexercise (userworkoutid, exerciseid, exercise_number) VALUES ($1, $2, $3) RETURNING id", workoutExercise.UserWorkoutId, workoutExercise.ExerciseId, workoutExercise.ExerciseNumber).Scan(&userWorkoutExerciseId)
	if err != nil {
		return 0, err
	}

	exerciseSetStmt, err := tx.Prepare("INSERT INTO userworkoutexerciseset (userworkoutexerciseid, reps, weight, is_done) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer exerciseSetStmt.Close()

	for _, set := range workoutExercise.UserWorkoutExerciseSets {
		_, err := exerciseSetStmt.Exec(userWorkoutExerciseId, set.Reps, set.Weight, set.IsDone)
		if err != nil {
			return 0, err
		}
	}

	return userWorkoutExerciseId, tx.Commit()
}

func (controller *Controller) DetermineExerciseNumberForNewUserWorkoutExercise(userWorkoutId int) (int, error) {
	var exerciseNumber int

	// When no exercise number is found return 0.
	err := controller.db.QueryRow(
		`SELECT COALESCE(MAX(exercise_number), 0)
			 FROM userworkoutexercise
			 WHERE userWorkoutId = $1`,
		userWorkoutId,
	).Scan(&exerciseNumber)

	if err != nil {
		return 0, err
	}

	return exerciseNumber + 1, err
}

func (controller *Controller) FindUserIdForUserWorkoutExerciseId(id int) (int, error) {
	rows, err := controller.db.Query("SELECT uw.userid FROM userworkoutexercise uwe JOIN userworkout uw ON uw.id = uwe.userworkoutid WHERE uwe.id = $1 LIMIT 1", id)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	userId := 0
	for rows.Next() {
		err := rows.Scan(&userId)
		if err != nil {
			return 0, err
		}
	}

	return userId, nil
}
