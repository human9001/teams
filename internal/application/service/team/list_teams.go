package team

import (
	"context"
	"log/slog"

	"github.com/human9001/teams/internal/domain/team"
)

func (s *TeamService) ListTeams(ctx context.Context, ownerID int64) ([]team.Team, error) {
	slog.Info("ListTeams service")
	return s.repo.ListByUser(ctx, ownerID)
}
