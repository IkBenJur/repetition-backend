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
	_, err := controller.db.Exec("INSERT INTO userWorkout (name) VALUES ($1)", workout.Name)
	if err != nil {
		return err
	}

	return nil
}