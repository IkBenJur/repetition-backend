package userWorkoutExerciseSet

import (
	"context"
	"time"

	base "github.com/IkBenJur/repetition-backend/service/baseController"
	types "github.com/IkBenJur/repetition-backend/types/userWorkout"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserWorkoutExerciseSetController struct {
	*base.BaseController[types.UserWorkoutExerciseSet]
}

var columnDefinitions = []base.ColumnDefinitionInterface{
	base.NewPrimaryKeyColumnDefinition(
		"id",
		false,
		func(w *types.UserWorkoutExerciseSet) **int64 { return &w.Id },
	),
	base.NewColumnDefinition(
		"user_id",
		false,
		func(w *types.UserWorkoutExerciseSet) **int64 { return &w.UserId },
	),
	base.NewColumnDefinition(
		"created_at",
		false,
		func(w *types.UserWorkoutExerciseSet) **time.Time { return &w.CreatedAt },
	),
	base.NewColumnDefinition(
		"updated_at",
		true,
		func(w *types.UserWorkoutExerciseSet) **time.Time { return &w.UpdatedAt },
	),
	base.NewColumnDefinition(
		"user_workout_exercise_id",
		true,
		func(w *types.UserWorkoutExerciseSet) *int { return &w.UserWorkoutExerciseId },
	),
	base.NewColumnDefinition(
		"reps",
		true,
		func(w *types.UserWorkoutExerciseSet) **int { return &w.Reps },
	),
	base.NewColumnDefinition(
		"set_number",
		true,
		func(w *types.UserWorkoutExerciseSet) **int { return &w.SetNumber },
	),
	base.NewColumnDefinition(
		"is_done",
		true,
		func(w *types.UserWorkoutExerciseSet) *bool { return &w.IsDone },
	),
	base.NewColumnDefinition(
		"load_prescription_id",
		true,
		func(w *types.UserWorkoutExerciseSet) *bool { return &w.IsDone },
	),
}

func NewUserWorkoutExerciseSetController(db *pgxpool.Pool) *UserWorkoutExerciseSetController {
	return &UserWorkoutExerciseSetController{
		BaseController: base.NewBaseController[types.UserWorkoutExerciseSet](db, "user_workout_exercise_set", columnDefinitions),
	}
}

func (controller *UserWorkoutExerciseSetController) Save(ctx context.Context, entity *types.UserWorkoutExerciseSet) (int64, error) {
	entity.Touch()

	if entity.IsNew() {
		return controller.Save(ctx, entity)
	} else {
		return controller.Update(ctx, entity)
	}
}

func (controller *UserWorkoutExerciseSetController) DetermineSetNumberForNewUserWorkoutExerciseSet(
	ctx context.Context,
	userWorkoutExerciseId int) (int, error) {
	var setNumber int

	// When no set number is found return 0.
	err := controller.DB.QueryRow(
		ctx,
		`SELECT COALESCE(MAX(set_number), 0)
			 FROM userworkoutexerciseset
			 WHERE userworkoutexerciseid = $1`,
		userWorkoutExerciseId,
	).Scan(&setNumber)

	if err != nil {
		return 0, err
	}

	return setNumber + 1, err
}
