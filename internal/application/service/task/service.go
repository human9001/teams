package task

import (
	"context"

	"github.com/human9001/teams/internal/domain/task/repository"
	"github.com/human9001/teams/internal/domain/team"
)

type ICache interface {
	Get(ctx context.Context, teamID int64, status, assigneeID string, page, limit int64, dst any) (bool, error)
	Set(ctx context.Context, teamID int64, status, assigneeID string, page, limit int64, value any) error
	BumpVersion(ctx context.Context, teamID int64) (int64, error)
}

type TaskService struct {
	repo   repository.IRepository
	member team.IMemberRepository
	cache  ICache
}

func NewTaskService(repo repository.IRepository, m team.IMemberRepository, c ICache) *TaskService {
	return &TaskService{
		repo:   repo,
		member: m,
		cache:  c,
	}
}
