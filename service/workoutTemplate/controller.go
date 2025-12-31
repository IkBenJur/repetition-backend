package workouttemplate

import (
	"database/sql"

	"github.com/IkBenJur/repetition-backend/types"
)

type Controller struct {
	db *sql.DB
}

func NewController(db *sql.DB) *Controller {
	return &Controller{
		db: db,
	}
}

func (controller *Controller) CreateNewTemplateWorkout(templateWorkout *types.TemplateWorkout) (int, error) {
	tx, err := controller.db.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	workoutTemplateStmt, err := newWorkoutTemplateStatement(tx)
	if err != nil {
		return -1, err
	}

	exerciseTemplateStmt, err := newTemplateExerciseStatement(tx)
	if err != nil {
		return -1, err
	}

	templateSetStmt, err := newTemplateSetStatement(tx)
	if err != nil {
		return -1, err
	}

	// Create new templateWorkout and set the new ID
	err = workoutTemplateStmt.QueryRow(templateWorkout.Name).Scan(templateWorkout.Id)
	if err != nil {
		return -1, err
	}

	for i, templateExercise := range templateWorkout.Exercises {
		err = exerciseTemplateStmt.
			QueryRow(
				templateExercise.ExerciseId,
				templateExercise.TemplateWorkoutId,
			).
			Scan(templateWorkout.Exercises[i].Id)
		if err != nil {
			return -1, err
		}

		for j, templateSet := range templateExercise.TemplateSets {
			templateSetStmt.
				QueryRow(
					templateSet.RepGoal,
					templateSet.WeightGoal,
				).Scan(templateWorkout.Exercises[i].TemplateSets[j].Id)
		}
	}

	return *templateWorkout.Id, nil
}

func newWorkoutTemplateStatement(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`INSERT INTO workout_template (name)
						VALUES ($1)
						RERTNING id`)
}

func newTemplateExerciseStatement(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`INSERT INTO template_workout_exercise (exercise_id, workout_template_id)
						VALUES ($1, $2)
						RETURNING id`)
}

func newTemplateSetStatement(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`INSERT INTO template_exercise_set (rep_goal, weight_goal)
						VALUES ($1, $2)
						RETURNING id`)
}
