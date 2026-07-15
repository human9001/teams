package team

import (
	"context"
	"log/slog"

	"github.com/human9001/teams/internal/domain/team"
)

func (s *TeamService) CreateTeam(ctx context.Context, name string, ownerID int64) (*team.Team, error) {
	slog.Info("CreateTeam Service")
	t, err := team.New(name, ownerID)
	if err != nil {
		return nil, err
	}
	teamId, err := s.repo.Create(ctx, t)
	if err != nil {
		return nil, err
	}
	t.SetID(teamId)
	return t, nil
}
