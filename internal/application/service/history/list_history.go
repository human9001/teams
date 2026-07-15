package taskhistory

import (
	"context"

	"github.com/human9001/teams/internal/application/service/history/dto"
	domainErrors "github.com/human9001/teams/internal/domain/errors"
)

func (s TaskHistoryService) GetTaskHistory(ctx context.Context, taskID int64) (dto.ListHistoryResult, error) {
	items, err := s.repo.ListByTaskID(ctx, taskID)
	if err != nil {
		return dto.ListHistoryResult{}, err
	}
	if len(items) == 0 {
		return dto.ListHistoryResult{}, domainErrors.ErrNotFound
	}
	return dto.ListHistoryResult{Items: items}, nil
}
