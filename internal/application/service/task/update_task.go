package task

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/human9001/teams/internal/application/service/task/input"
	domainErrs "github.com/human9001/teams/internal/domain/errors"
	"github.com/human9001/teams/pkg/helpers"
)

func (s *TaskService) UpdateTask(ctx context.Context, in input.UpdateTaskInput) error {
	if in.TaskID == 0 || in.UserID == 0 {
		return domainErrs.ErrInvalidInput
	}

	t, err := s.repo.ByID(ctx, in.TaskID)
	if err != nil {
		return err
	}

	ok, err := s.member.IsTeamMember(ctx, t.TeamID(), in.UserID)
	if err != nil {
		return err
	}
	if !ok {
		return domainErrs.ErrForbidden
	}

	if in.Title != nil {
		v := strings.TrimSpace(*in.Title)
		if v == "" {
			return errors.New("title cannot be empty")
		}
		t.SetTitle(v)
	}
	if in.Description != nil {
		t.SetDescription(strings.TrimSpace(*in.Description))
	}
	t.SetTitle(*helpers.PtrOrNIl(in.Title))
	t.SetDescription(*helpers.PtrOrNIl(new(strings.TrimSpace(*in.Description))))
	t.SetStatus(*helpers.PtrOrNIl(in.Status))
	t.SetPriority(*helpers.PtrOrNIl(in.Priority))
	t.SetAssigneeID(helpers.PtrOrNIl(in.AssigneeID))

	_, err = s.cache.BumpVersion(ctx, t.TeamID())
	if err != nil {
		slog.Info("bumpVersion", "error", err)
	}

	return s.repo.Update(ctx, t)
}
