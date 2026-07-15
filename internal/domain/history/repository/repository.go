package repository

import (
	"context"

	"github.com/human9001/teams/internal/domain/history/model"
)

type IRepository interface {
	ListByTaskID(ctx context.Context, taskID int64) ([]model.HistoryItem, error)
}
