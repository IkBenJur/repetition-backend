package types

import (
	"errors"
	"time"

	"github.com/IkBenJur/repetition-backend/types"
)

type NewUserWorkoutPayload struct {
	types.BasePayloadEntity
	Name                 string                       `json:"name" validate:"required"`
	DateStart            *time.Time                   `json:"dateStart"`
	DateEnd              *time.Time                   `json:"dateEnd"`
	UserWorkoutExercises []UserWorkoutExercisePayload `json:"userWorkoutExercises"`
}

func (payload *NewUserWorkoutPayload) ToEntity() (*UserWorkout, error) {
	exercises := make([]*UserWorkoutExercise, len(payload.UserWorkoutExercises))
	for i, exerciseEntity := range payload.UserWorkoutExercises {
		exercise, err := exerciseEntity.ToEntity()
		if err != nil {
			return nil, err
		}
		exercises[i] = exercise
	}

	baseEntity := payload.ToBaseEntity()
	if baseEntity == nil {
		return nil, errors.New("base entity is nil")
	}

	return &UserWorkout{
		BaseEntity:           *baseEntity,
		Name:                 payload.Name,
		DateStart:            payload.DateStart,
		DateEnd:              payload.DateEnd,
		UserWorkoutExercises: exercises,
	}, nil
}

type UserWorkout struct {
	types.BaseEntity
	Name                 string
	DateStart            *time.Time
	DateEnd              *time.Time
	UserWorkoutExercises []*UserWorkoutExercise
}

func (uw *UserWorkout) MarkStarted() {
	if uw.DateStart != nil {
		return
	}
	dateStart := time.Now()
	uw.DateStart = &dateStart
}

func (uw *UserWorkout) MarkFinished() {
	if uw.DateEnd != nil {
		return
	}
	dateEnd := time.Now()
	uw.DateEnd = &dateEnd
}
