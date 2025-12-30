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
	ID             int
	UserWorkoutId  int
	ExerciseId     int
	ExerciseNumber *int
	CreatedAt      time.Time

	// Joined fields
	UserWorkoutExerciseSets []*UserWorkoutExerciseSet
	ExerciseName            *string
}

type UserWorkoutExerciseSet struct {
	ID                    int
	UserWorkoutExerciseId int
	Reps                  *int
	Weight                *float64
	SetNumber             *int
	IsDone                bool
	CreatedAt             time.Time
}

type UserWorkoutExercisePayload struct {
	ExerciseId              int                             `json:"exerciseId" validate:"required"`
	ExerciseName            *string                         `json:"exerciseName"`
	UserWorkoutId           int                             `json:"userWorkoutId" validate:"required"`
	ExerciseNumber          *int                            `json:"exerciseNumber"`
	UserWorkoutExerciseSets []UserWorkoutExerciseSetPayload `json:"userWorkoutExerciseSets"`
}

func (payload UserWorkoutExercisePayload) ToEntity() *UserWorkoutExercise {
	sets := make([]*UserWorkoutExerciseSet, len(payload.UserWorkoutExerciseSets))
	for i, set := range payload.UserWorkoutExerciseSets {
		sets[i] = set.ToEntity()
	}

	return &UserWorkoutExercise{
		UserWorkoutExerciseSets: sets,
		ExerciseId:              payload.ExerciseId,
		ExerciseName:            payload.ExerciseName,
		UserWorkoutId:           payload.UserWorkoutId,
		ExerciseNumber:          payload.ExerciseNumber,
	}
}

type UserWorkoutExerciseSetPayload struct {
	ID                    *int     `json:"id"`
	UserWorkoutExerciseId int      `json:"userWorkoutExerciseId" validate:"required"`
	Reps                  *int     `json:"reps"`
	Weight                *float64 `json:"weight"`
	SetNumber             *int     `json:"setNumber"`
	IsDone                *bool    `json:"isDone"`
}

func (payload UserWorkoutExerciseSetPayload) ToEntity() *UserWorkoutExerciseSet {
	id := 0
	if payload.ID != nil {
		id = *payload.ID
	}

	isDone := false
	if payload.IsDone != nil {
		isDone = *payload.IsDone
	}

	return &UserWorkoutExerciseSet{
		ID:                    id,
		UserWorkoutExerciseId: payload.UserWorkoutExerciseId,
		Reps:                  payload.Reps,
		Weight:                payload.Weight,
		SetNumber:             payload.SetNumber,
		IsDone:                isDone,
	}
}

func (payload UserWorkoutExerciseSetPayload) IsUpdate() bool {
	return payload.ID != nil && *payload.ID > 0
}

type NewUserWorkoutPayload struct {
	Name                 string                       `json:"name" validate:"required"`
	UserId               int                          `json:"userId"`
	UserWorkoutExercises []UserWorkoutExercisePayload `json:"userWorkoutExercises"`
}

func (payload *NewUserWorkoutPayload) ToEntity() *UserWorkout {
	exercises := make([]*UserWorkoutExercise, len(payload.UserWorkoutExercises))
	for i, exercise := range payload.UserWorkoutExercises {
		exercises[i] = exercise.ToEntity()
	}

	return &UserWorkout{
		Name:                 payload.Name,
		UserId:               payload.UserId,
		UserWorkoutExercises: exercises,
	}
}
