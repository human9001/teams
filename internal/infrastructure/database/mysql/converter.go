package mysql

import (
	"time"

	"github.com/human9001/teams/internal/domain/team"
)

type teamRow struct {
	ID        int64
	Name      string
	OwnerID   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func toDomainTeam(r teamRow) team.Team {
	return team.NewFromPersistence(r.ID, r.Name, r.OwnerID, r.CreatedAt, r.UpdatedAt)
}
