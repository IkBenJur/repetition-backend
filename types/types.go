package types

import (
	"time"
)

type UserController interface {
	GetUserByUsername(username string) (*User, error)
	SaveUser(user User) error
	GetUserById(id int) (*User, error)
}

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	ID int
	Username string
	Password string
	CreatedAt time.Time
}