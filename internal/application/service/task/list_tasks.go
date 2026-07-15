package task

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/human9001/teams/internal/application/converter"
	"github.com/human9001/teams/internal/application/service/task/input"
)

func (s *TaskService) ListTasks(ctx context.Context, in input.ListTasksInput) (input.ListTasksResult, error) {
	var cached input.ListTasksResult
	var st string
	var aId string

	if in.Status != nil {
		st = in.Status.String()
	}
	if in.AssigneeID != nil {
		aId = fmt.Sprintf("%v", in.AssigneeID)
	}
	ok, err := s.cache.Get(ctx, in.TeamID, st, aId, in.Page, in.Limit, &cached)
	if err == nil && ok {
		return cached, nil
	}
	if err != nil {
		slog.Error("tasks from cache", "error", err)
	}

	res, total, err := s.repo.List(ctx, in)
	if err != nil {
		return input.ListTasksResult{}, err
	}

	listResult := input.ListTasksResult{
		Items: converter.FromModelToListTasksItem(res),
		Total: total,
		Page:  in.Page,
		Limit: in.Limit,
	}

	err = s.cache.Set(ctx, in.TeamID, st, aId, in.Page, in.Limit, listResult)
	if err != nil {
		slog.Error("list task set cache", "error", "err")
	}

	return listResult, nil
}
