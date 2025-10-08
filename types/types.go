package types

import (
	"time"
)

type UserController interface {
	GetUserByUsername(username string) (*User, error)
	SaveUser(user User) error
}

type RegisterUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID int
	Username string
	Password string
	CreatedAt time.Time
}