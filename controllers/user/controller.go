package controller

import (
	"database/sql"

	"github.com/IkBenJur/repetition-backend/types"
)

type Controller struct {
	db *sql.DB
}

func NewController(db *sql.DB) *Controller {
	return &Controller{ db: db }
}

func (controller *Controller) GetUserByUsername(username string) (*types.User, error) {
	rows, err := controller.db.Query("SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		return nil, err
	}

	user := new(types.User)

	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (controller *Controller) SaveUser(user types.User) error {
	return nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}