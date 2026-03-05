package types

import (
	"errors"

	"github.com/IkBenJur/repetition-backend/types"
)

type ExercisePayload struct {
	types.BasePayloadEntity
	Name        string `json:"name" validate:"required"`
	MuscleGroup string `json:"muscleGroup" validate:"required"`
}

func (payload *ExercisePayload) ToEntity() (*Exercise, error) {
	baseEntity := payload.ToBaseEntity()
	if baseEntity == nil {
		return nil, errors.New("base entity is nil")
	}

	return &Exercise{
		BaseEntity:  *baseEntity,
		Name:        payload.Name,
		MuscleGroup: payload.MuscleGroup,
	}, nil
}

type Exercise struct {
	types.BaseEntity
	Name        string `json:"name"`
	MuscleGroup string `json:"muscleGroup"`
}

