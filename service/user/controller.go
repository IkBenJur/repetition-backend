package user

import (
	"context"
	"fmt"
	"time"

	base "github.com/IkBenJur/repetition-backend/service/baseController"
	types "github.com/IkBenJur/repetition-backend/types/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

var columnDefinitions = []base.ColumnDefinitionInterface{
	base.NewPrimaryKeyColumnDefinition(
		"id",
		false,
		func(w *types.User) **int64 { return &w.Id },
	),
	base.NewColumnDefinition(
		"created_at",
		false,
		func(w *types.User) **time.Time { return &w.CreatedAt },
	),
	base.NewColumnDefinition(
		"updated_at",
		true,
		func(w *types.User) **time.Time { return &w.UpdatedAt },
	),
	base.NewColumnDefinition(
		"username",
		true,
		func(w *types.User) *string { return &w.Username },
	),
	base.NewColumnDefinition(
		"password",
		true,
		func(w *types.User) *string { return &w.Password },
	),
	base.NewColumnDefinition(
		"active_userworkout_id",
		true,
		func(w *types.User) **int64 { return &w.ActiveUserWorkoutId },
	),
}

type UserController struct {
	*base.BaseController[types.User]
}

func NewUserController(db *pgxpool.Pool) *UserController {
	return &UserController{
		BaseController: base.NewBaseController[types.User](db, "user_table", columnDefinitions),
	}
}

func (controller *UserController) FindById(ctx context.Context, id int64) (*types.User, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		controller.BaseController.SelectColumns,
		controller.BaseController.TableName,
	)

	rows, err := controller.BaseController.DB.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return controller.BaseController.ScanRow(rows)
}

func (controller *UserController) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE username = $1",
		controller.BaseController.SelectColumns,
		controller.BaseController.TableName,
	)

	rows, err := controller.BaseController.DB.Query(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return controller.BaseController.ScanRow(rows)
}

func (controller *UserController) Save(ctx context.Context, user *types.User) (int64, error) {
	user.Touch()

	if user.IsNew() {
		return controller.Create(ctx, user)
	} else {
		return controller.Update(ctx, user)
	}
}
