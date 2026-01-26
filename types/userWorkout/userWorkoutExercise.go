package types

import (
	"errors"

	"github.com/IkBenJur/repetition-backend/types"
)

type UserWorkoutExercisePayload struct {
	types.BasePayloadEntity
	ExerciseId              int                             `json:"exerciseId" validate:"required"`
	ExerciseName            *string                         `json:"exerciseName"`
	UserWorkoutId           int                             `json:"userWorkoutId" validate:"required"`
	ExerciseNumber          *int                            `json:"exerciseNumber"`
	UserWorkoutExerciseSets []UserWorkoutExerciseSetPayload `json:"userWorkoutExerciseSets"`
}

func (payload *UserWorkoutExercisePayload) ToEntity() (*UserWorkoutExercise, error) {
	sets := make([]*UserWorkoutExerciseSet, len(payload.UserWorkoutExerciseSets))
	for i, setEntity := range payload.UserWorkoutExerciseSets {

		set, err := setEntity.ToEntity()
		if err != nil {
			return nil, err
		}

		sets[i] = set
	}

	baseEntity := payload.ToBaseEntity()
	if baseEntity == nil {
		return nil, errors.New("base entity is nil")
	}

	return &UserWorkoutExercise{
		BaseEntity:              *baseEntity,
		UserWorkoutId:           payload.UserWorkoutId,
		ExerciseId:              payload.ExerciseId,
		ExerciseNumber:          payload.ExerciseNumber,
		ExerciseName:            payload.ExerciseName,
		UserWorkoutExerciseSets: sets,
	}, nil
}

type UserWorkoutExercise struct {
	types.BaseEntity
	UserWorkoutId  int
	ExerciseId     int
	ExerciseNumber *int

	// Joined fields
	UserWorkoutExerciseSets []*UserWorkoutExerciseSet
	ExerciseName            *string
}
