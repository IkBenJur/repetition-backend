package workouttemplate

import (
	"database/sql"
	"fmt"

	loadprescription "github.com/IkBenJur/repetition-backend/service/loadPrescription"
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

	loadPrescriptionController := loadprescription.NewController(controller.db)

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

	loadPrescriptionStmt, err := loadPrescriptionController.CreateLoadPrescriptionStatement(tx)
	if err != nil {
		return -1, err
	}

	fixLoadPrescriptionStmt, err := loadPrescriptionController.CreateFixedLoadPrescriptionStatement(tx)
	if err != nil {
		return -1, err
	}

	// Create new templateWorkout and set the new ID
	err = workoutTemplateStmt.
		QueryRow(
			templateWorkout.Name,
			templateWorkout.UserId,
		).
		Scan(&templateWorkout.Id)
	if err != nil {
		return -1, err
	}

	for i, templateExercise := range templateWorkout.Exercises {

		templateExercise.TemplateWorkoutId = &templateWorkout.Id

		err = exerciseTemplateStmt.
			QueryRow(
				templateExercise.ExerciseId,
				templateExercise.TemplateWorkoutId,
			).
			Scan(&templateWorkout.Exercises[i].Id)
		if err != nil {
			return -1, err
		}

		for j, templateSet := range templateExercise.TemplateSets {

			if !loadPrescriptionController.IsValidLoadPrescriptionType(templateSet.LoadPresciptionType) {
				return -1, fmt.Errorf("Invalid load type ID found: %v", templateSet.LoadPresciptionType)
			}

			// Create new loadPrescription for set
			err := loadPrescriptionStmt.
				QueryRow(templateSet.LoadPresciptionType).
				Scan(&templateSet.PrescriptionId)
			if err != nil {
				return -1, err
			}

			// Add to correct specified table
			if *templateSet.LoadPresciptionType == types.FIXED {

				templateSet.FixedLoadPrescription.Id = &templateSet.PrescriptionId

				_, err := fixLoadPrescriptionStmt.
					Exec(
						templateSet.PrescriptionId,
						templateSet.FixedLoadPrescription.Weight,
					)
				if err != nil {
					return -1, err
				}
			} // TODO Same if statements for OneRepMax and RPE types

			err = templateSetStmt.
				QueryRow(
					templateSet.RepGoal,
					templateSet.PrescriptionId,
				).Scan(&templateWorkout.Exercises[i].TemplateSets[j].Id)
			if err != nil {
				return -1, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return templateWorkout.Id, nil
}

func newWorkoutTemplateStatement(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`INSERT INTO workout_template (name, user_id)
						VALUES ($1, $2)
						RETURNING id`)
}

func newTemplateExerciseStatement(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`INSERT INTO template_workout_exercise (exercise_id, workout_template_id)
						VALUES ($1, $2)
						RETURNING id`)
}

func newTemplateSetStatement(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`INSERT INTO template_exercise_set (rep_goal, load_prescription_id)
						VALUES ($1, $2)
						RETURNING id`)
}
