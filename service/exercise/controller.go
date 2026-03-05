package exercise

import (
	"context"
	"time"

	base "github.com/IkBenJur/repetition-backend/service/baseController"
	"github.com/IkBenJur/repetition-backend/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

var columnDefinitions = []base.ColumnDefinitionInterface{
	base.NewPrimaryKeyColumnDefinition(
		"id",
		false,
		func(e *types.Exercise) **int64 { return &e.Id },
	),
	base.NewColumnDefinition(
		"created_at",
		false,
		func(e *types.Exercise) **time.Time { return &e.CreatedAt },
	),
	base.NewColumnDefinition(
		"updated_at",
		true,
		func(e *types.Exercise) **time.Time { return &e.UpdatedAt },
	),
	base.NewColumnDefinition(
		"name",
		true,
		func(e *types.Exercise) *string { return &e.Name },
	),
	base.NewColumnDefinition(
		"muscle_group",
		true,
		func(e *types.Exercise) *string { return &e.MuscleGroup },
	),
}

type ExerciseController struct {
	*base.BaseController[types.Exercise]
}

func NewExerciseController(db *pgxpool.Pool) *ExerciseController {
	return &ExerciseController{
		BaseController: base.NewBaseController[types.Exercise](db, "exercise", columnDefinitions),
	}
}

func (c *ExerciseController) Save(ctx context.Context, entity *types.Exercise) (int64, error) {
	entity.Touch()

	if entity.IsNew() {
		return c.Create(ctx, entity)
	}

	return c.Update(ctx, entity)
}

func (c *ExerciseController) FindAll(ctx context.Context) ([]*types.Exercise, error) {
	query := "SELECT " + c.SelectColumns + " FROM " + c.TableName

	rows, err := c.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return c.ScanRows(rows)
}

func (c *ExerciseController) FindById(ctx context.Context, id int64) (*types.Exercise, error) {
	query := "SELECT " + c.SelectColumns + " FROM " + c.TableName + " WHERE id = $1"

	rows, err := c.DB.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return c.ScanRow(rows)
}

