package exercise

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

func (c *Controller) GetAllExercise() ([]types.Exercise, error) {
	rows, err := c.db.Query("SELECT * FROM exercise")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exercises := make([]types.Exercise, 0)

	for rows.Next() {
		exercise, err := scanRowIntoExercise(rows)
		if err != nil {
			return nil, err
		}

		exercises = append(exercises, *exercise)
	}

	return exercises, nil
}

func (c *Controller) GetExerciseById(id int) (*types.Exercise, error) {
	rows, err := c.db.Query("SELECT * FROM exercise WHERE id = $1 LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exercise := new(types.Exercise)

	for rows.Next() {
		exercise, err = scanRowIntoExercise(rows)
		if err != nil {
			return nil, err
		}
	}

	return exercise, nil
}

func (c *Controller) SaveExercise(exercise types.Exercise) error {
	_, err := c.db.Exec("INSERT INTO exercise (name, muscleGroup) VALUES ($1, $2)", exercise.Name, exercise.MuscleGroup)
	if err != nil {
		return err
	}

	return nil
}

func scanRowIntoExercise(rows *sql.Rows) (*types.Exercise, error) {
	exercise := new(types.Exercise)

	err := rows.Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.MuscleGroup,
		&exercise.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return exercise, nil
}
