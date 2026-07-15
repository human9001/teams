package taskhistory

import (
	"github.com/human9001/teams/internal/domain/history/repository"
)

type TaskHistoryService struct {
	repo repository.IRepository
}

func NewTaskHistoryService(repo repository.IRepository) *TaskHistoryService {
	return &TaskHistoryService{
		repo: repo,
	}
}
