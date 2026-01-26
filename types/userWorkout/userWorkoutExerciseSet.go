package types

import (
	"errors"

	"github.com/IkBenJur/repetition-backend/types"
)

type UserWorkoutExerciseSetPayload struct {
	types.BasePayloadEntity
	UserWorkoutExerciseId int   `json:"userWorkoutId"  validate:"required"`
	Reps                  *int  `json:"reps"`
	SetNumber             *int  `json:"setNumber"`
	IsDone                *bool `json:"isDone"`
	// TODO Implement the loadPrescriptions
}

func (payload *UserWorkoutExerciseSetPayload) ToEntity() (*UserWorkoutExerciseSet, error) {

	baseEntity := payload.ToBaseEntity()
	if baseEntity == nil {
		return nil, errors.New("base entity is nil")
	}

	isDone := false
	if payload.IsDone != nil {
		isDone = *payload.IsDone
	}

	return &UserWorkoutExerciseSet{
		BaseEntity:            *baseEntity,
		UserWorkoutExerciseId: payload.UserWorkoutExerciseId,
		Reps:                  payload.Reps,
		SetNumber:             payload.SetNumber,
		IsDone:                isDone,
	}, nil
}

type UserWorkoutExerciseSet struct {
	types.BaseEntity
	UserWorkoutExerciseId int
	Reps                  *int
	SetNumber             *int
	IsDone                bool
}
