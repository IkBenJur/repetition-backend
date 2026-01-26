package types

import (
	"errors"

	"github.com/IkBenJur/repetition-backend/types"
)

type NewUserWorkoutPayload struct {
	types.BasePayloadEntity
	Name                 string                       `json:"name" validate:"required"`
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
		UserWorkoutExercises: exercises,
	}, nil
}

type UserWorkout struct {
	types.BaseEntity
	Name                 string
	UserWorkoutExercises []*UserWorkoutExercise
}
