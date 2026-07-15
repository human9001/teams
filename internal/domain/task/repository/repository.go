package repository

import (
	"context"

	"github.com/human9001/teams/internal/application/service/task/input"
	"github.com/human9001/teams/internal/domain/task/model"
)

type IRepository interface {
	Save(ctx context.Context, t *model.Task) (int64, error)
	ByID(ctx context.Context, id int64) (model.Task, error)
	List(ctx context.Context, filter input.ListTasksInput) ([]model.Task, int64, error)
	Update(ctx context.Context, t model.Task) error
}
