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
	err := controller.db.QueryRow("INSERT INTO userworkoutexerciseset (userworkoutexerciseid, reps, weight, set_number, is_done) VALUES ($1, $2, $3, $4, $5) RETURNING id", workoutExerciseSet.UserWorkoutExerciseId, workoutExerciseSet.Reps, workoutExerciseSet.Weight, workoutExerciseSet.SetNumber, workoutExerciseSet.IsDone).Scan(&workoutExerciseSetId)

	return workoutExerciseSetId, err
}

func (controller *Controller) DetermineSetNumberForNewUserWorkoutExerciseSet(userWorkoutExerciseId int) (int, error) {
	var setNumber int

	// When no set number is found return 0.
	err := controller.db.QueryRow(
		`SELECT COALESCE(MAX(set_number), 0)
			 FROM userworkoutexerciseset
			 WHERE userworkoutexerciseid = $1`,
		userWorkoutExerciseId,
	).Scan(&setNumber)

	if err != nil {
		return 0, err
	}

	return setNumber + 1, err
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
	_, err := controller.db.Exec("UPDATE userworkoutexerciseset SET reps = $1, weight = $2, is_done = $4 WHERE id = $3", workoutExerciseSet.Reps, workoutExerciseSet.Weight, workoutExerciseSet.ID, workoutExerciseSet.IsDone)
	return err
}
