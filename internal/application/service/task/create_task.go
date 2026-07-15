package task

import (
	"context"
	"errors"
	"log/slog"

	"github.com/human9001/teams/internal/application/service/task/input"
	"github.com/human9001/teams/internal/domain/task/model"
)

func (s *TaskService) CreateTask(ctx context.Context, userId int64, in input.CreateTaskInput) (*model.Task, error) {
	slog.Info("Creating task", "userId", userId, "title", in.Title, "description", in.Description, "priority", in.Priority)
	ok, err := s.member.IsTeamMember(ctx, in.TeamID, userId)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("forbidden: not a team member")
	}

	t, err := model.New(in.TeamID, userId, in.Title, in.Description, in.Priority)
	if err != nil {
		return nil, err
	}

	taskId, err := s.repo.Save(ctx, t)
	if err != nil {
		return nil, err
	}
	t.SetID(taskId)
	_, err = s.cache.BumpVersion(ctx, t.TeamID())
	if err != nil {
		slog.Info("bumpVersion", "error", err)
	}
	return t, nil
}
