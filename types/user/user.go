package types

import (
	"time"
)

type UserPayload struct {
	Id                  *int64     `json:"id"`
	ActiveUserWorkoutId *int64     `json:"active_user_workout_id"`
	Username            string     `json:"username" validate:"required"`
	Password            string     `json:"password" validate:"required,min=3,max=130"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
}

func (payload *UserPayload) IsUpdate() bool {
	return payload.Id != nil && *payload.Id > 0
}

func (payload *UserPayload) IsNew() bool {
	return payload.Id == nil
}

func (payload *UserPayload) ToEntity() *User {
	user := &User{}

	if payload.Id != nil {
		user.Id = payload.Id
	}

	if payload.ActiveUserWorkoutId != nil {
		user.ActiveUserWorkoutId = payload.ActiveUserWorkoutId
	}

	if payload.CreatedAt != nil {
		user.CreatedAt = payload.CreatedAt
	}

	if payload.UpdatedAt != nil {
		user.UpdatedAt = payload.UpdatedAt
	}

	user.Username = payload.Username
	user.Password = payload.Password

	return user
}

type User struct {
	Id                  *int64
	ActiveUserWorkoutId *int64
	Username            string
	Password            string
	CreatedAt           *time.Time
	UpdatedAt           *time.Time
}

func (b *User) IsNew() bool {
	return b.Id == nil
}

func (b *User) IsUpdate() bool {
	return b.Id != nil && *b.Id > 0
}

func (b *User) MarkCreated() {
	now := time.Now()
	b.CreatedAt = &now
	b.UpdatedAt = &now
}

func (b *User) MarkUpdated() {
	now := time.Now()
	b.UpdatedAt = &now
}

// This function can be used to correctly update the right timestamp
func (b *User) Touch() {
	if b.IsNew() {
		b.MarkCreated()
	} else {
		b.MarkUpdated()
	}
}
