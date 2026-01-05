package types

import "time"

type TemplateWorkout struct {
	Id        int
	UserId    int
	Name      string
	CreatedAt time.Time

	// Joined fields
	Exercises []*TemplateWorkoutExercise
}

type TemplateWorkoutExercise struct {
	Id                *int
	ExerciseId        int
	TemplateWorkoutId *int
	CreatedAt         time.Time

	// Joined fields
	TemplateSets []*TemplateExerciseSet
}

type TemplateExerciseSet struct {
	Id             int
	RepGoal        int
	PrescriptionId int

	// Joined fields
	LoadPresciptionType                 *LoadPresciptionType
	FixedLoadPrescription               *FixedLoadPrescription
	PercentageOneRepMaxLoadPrescription *PercentageOneRepMaxLoadPrescription
	RPELoadPrescription                 *RPELoadPrescription
}

type TemplateWorkoutPayload struct {
	Name string `json:"name" validate:"required"`

	TemplateExercises []TemplateWorkoutExercisePayload `json:"templateExercises"`
}

func (payload *TemplateWorkoutPayload) ToEntity() *TemplateWorkout {
	exercises := make([]*TemplateWorkoutExercise, len(payload.TemplateExercises))
	for i, exercise := range payload.TemplateExercises {
		exercises[i] = exercise.ToEntity()
	}

	return &TemplateWorkout{
		Name: payload.Name,

		Exercises: exercises,
	}
}

type TemplateWorkoutExercisePayload struct {
	ExerciseId int `json:"exerciseId" validate:"required"`

	TemplateSets []TemplateExerciseSetPayload `json:"templateSets"`
}

func (payload *TemplateWorkoutExercisePayload) ToEntity() *TemplateWorkoutExercise {
	sets := make([]*TemplateExerciseSet, len(payload.TemplateSets))
	for i, set := range payload.TemplateSets {
		sets[i] = set.ToEntity()
	}

	return &TemplateWorkoutExercise{
		ExerciseId:   payload.ExerciseId,
		TemplateSets: sets,
	}
}

type TemplateExerciseSetPayload struct {
	RepGoal               int                               `json:"repGoal" validate:"required"`
	LoadPresciptionTypeId LoadPresciptionType               `json:"loadPrescriptionTypeId" validate:"required"`
	FixedLoadPrescription *FixedLoadPrescriptionTypePayload `json:"fixedLoadPrescription" validate:"required"`
}

func (payload *TemplateExerciseSetPayload) ToEntity() *TemplateExerciseSet {
	return &TemplateExerciseSet{
		RepGoal:             payload.RepGoal,
		LoadPresciptionType: &payload.LoadPresciptionTypeId,

		// Fields may be nill
		FixedLoadPrescription: payload.FixedLoadPrescription.ToEntity(),
	}
}

type FixedLoadPrescriptionTypePayload struct {
	Weight float64 `json:"weight"`
}

func (payload *FixedLoadPrescriptionTypePayload) ToEntity() *FixedLoadPrescription {
	if payload == nil {
		return nil
	}

	return &FixedLoadPrescription{
		Weight: payload.Weight,
	}
}

type PercentageOneRepMaxLoadPrescriptionPayload struct {
	Percentage float64
}

func (payload *PercentageOneRepMaxLoadPrescriptionPayload) ToEntity() *PercentageOneRepMaxLoadPrescription {
	if payload == nil {
		return nil
	}

	return &PercentageOneRepMaxLoadPrescription{
		Percentage: payload.Percentage,
	}
}
