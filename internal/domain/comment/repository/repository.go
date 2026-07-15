package repository

import (
	"context"

	"github.com/human9001/teams/internal/domain/comment/model"
)

type IRepository interface {
	ListByTaskID(ctx context.Context, taskID, limit, offset int64) ([]model.Comment, int64, error)
	UserHasAccess(ctx context.Context, taskID, userID int64) (bool, error)
}
