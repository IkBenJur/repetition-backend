package userWorkoutExercise

import (
	"context"
	"time"

	base "github.com/IkBenJur/repetition-backend/service/baseController"
	types "github.com/IkBenJur/repetition-backend/types/userWorkout"
	"github.com/jackc/pgx/v5/pgxpool"
)

var columnDefinitions = []base.ColumnDefinitionInterface{
	base.NewPrimaryKeyColumnDefinition(
		"id",
		false,
		func(w *types.UserWorkoutExercise) **int64 { return &w.Id },
	),
	base.NewColumnDefinition(
		"user_id",
		false,
		func(w *types.UserWorkoutExercise) **int64 { return &w.UserId },
	),
	base.NewColumnDefinition(
		"created_at",
		false,
		func(w *types.UserWorkoutExercise) **time.Time { return &w.CreatedAt },
	),
	base.NewColumnDefinition(
		"updated_at",
		true,
		func(w *types.UserWorkoutExercise) **time.Time { return &w.UpdatedAt },
	),
	base.NewColumnDefinition(
		"user_workout_id",
		true,
		func(w *types.UserWorkoutExercise) *int { return &w.UserWorkoutId },
	),
	base.NewColumnDefinition(
		"exercise_id",
		true,
		func(w *types.UserWorkoutExercise) *int { return &w.ExerciseId },
	),
	base.NewColumnDefinition(
		"exercise_number",
		true,
		func(w *types.UserWorkoutExercise) **int { return &w.ExerciseNumber },
	),

	// TODO Add joined column definition for exerciseName
}

type UserWorkoutExerciseController struct {
	*base.BaseController[types.UserWorkoutExercise]
}

func NewUserWorkoutExerciseController(db *pgxpool.Pool) *UserWorkoutExerciseController {
	return &UserWorkoutExerciseController{
		BaseController: base.NewBaseController[types.UserWorkoutExercise](db, "user_workout_exercise", columnDefinitions),
	}
}

func (controller *UserWorkoutExerciseController) Save(ctx context.Context, entity *types.UserWorkoutExercise) (int64, error) {
	entity.Touch()

	if entity.IsNew() {
		return controller.Create(ctx, entity)
	} else {
		return controller.Update(ctx, entity)
	}
}

func (controller *UserWorkoutExerciseController) DetermineExerciseNumberForNewUserWorkoutExercise(
	ctx context.Context,
	userWorkoutId int) (int, error) {
	var exerciseNumber int

	// When no exercise number is found return 0.
	err := controller.DB.QueryRow(
		ctx,
		`SELECT COALESCE(MAX(exercise_number), 0)
			 FROM userworkoutexercise
			 WHERE userWorkoutId = $1`,
		userWorkoutId,
	).Scan(&exerciseNumber)

	if err != nil {
		return 0, err
	}

	return exerciseNumber + 1, err
}

func (controller *UserWorkoutExerciseController) FindUserIdForUserWorkoutExerciseId(ctx context.Context, id int) (int, error) {
	rows, err := controller.DB.Query(
		ctx,
		"SELECT uw.userid FROM userworkoutexercise uwe JOIN userworkout uw ON uw.id = uwe.userworkoutid WHERE uwe.id = $1 LIMIT 1",
		id)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	userId := 0
	for rows.Next() {
		err := rows.Scan(&userId)
		if err != nil {
			return 0, err
		}
	}

	return userId, nil
}
