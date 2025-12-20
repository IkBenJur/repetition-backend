package types

import (
	"time"
)

type UserController interface {
	GetUserByUsername(username string) (*User, error)
	CreateNewUser(user User) error
	UpdateUser(user User) error
	GetUserById(id int) (*User, error)
	UpdateActiveUserWorkoutForUserId(userId int, activeUserWorkoutId int) error
}

type ExerciseController interface {
	GetAllExercise() ([]Exercise, error)
	GetExerciseById(id int) (*Exercise, error)
	SaveExercise(exercise Exercise) error
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
	ID                  int
	ActiveUserWorkoutId *int
	Username            string
	Password            string
	CreatedAt           time.Time
}

type NewExercisePayload struct {
	Name        string `json:"name" validate:"required"`
	MuscleGroup string `json:"muscleGroup" validate:"required"`
}

type Exercise struct {
	ID          int
	Name        string
	MuscleGroup string
	CreatedAt   time.Time
}

type UserWorkout struct {
	ID        int
	UserId    int
	Name      string
	DateStart time.Time
	DateEnd   *time.Time
	CreatedAt time.Time

	// Joined fields
	UserWorkoutExercises []*UserWorkoutExercise
}

type UserWorkoutExercise struct {
	ID            int
	UserWorkoutId int
	ExerciseId    int
	CreatedAt     time.Time

	// Joined fields
	UserWorkoutExerciseSets []*UserWorkoutExerciseSet
}

type UserWorkoutExerciseSet struct {
	ID                    int
	UserWorkoutExerciseId int
	Reps                  int
	Weight                float32
	CreatedAt             time.Time
}

type UserWorkoutExercisePayload struct {
	ExerciseId              int                             `json:"exerciseId" validate:"required"`
	UserWorkoutId           int                             `json:"userWorkoutId" validate:"required"`
	UserWorkoutExerciseSets []UserWorkoutExerciseSetPayload `json:"userWorkoutExerciseSets"`
}

type UserWorkoutExerciseSetPayload struct {
	UserWorkoutExerciseId int     `json:"userWorkoutExerciseId" validate:"required"`
	Reps                  int     `json:"reps" validate:"required"`
	Weight                float32 `json:"weight" validate:"required"`
}

type NewUserWorkoutPayload struct {
	Name                 string                       `json:"name"`
	UserId               int                          `json:"userId" validate:"required"`
	UserWorkoutExercises []UserWorkoutExercisePayload `json:"userWorkoutExercises"`
}
