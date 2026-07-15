package mysql

import (
	"context"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/human9001/teams/internal/domain/team"
)

type TeamRepository struct {
	db *sqlx.DB
	tx TxManager
}

func NewTeamRepository(db *sqlx.DB, tx TxManager) *TeamRepository {
	return &TeamRepository{db: db, tx: tx}
}

func (r *TeamRepository) Create(ctx context.Context, t *team.Team) (int64, error) {
	slog.Info("Create Team Repo")
	var teamId int64
	err := r.tx.Do(ctx, func(txCtx context.Context) error {
		res, err := r.db.ExecContext(ctx, `
		INSERT INTO teams (name, owner_id, created_at)
		VALUES (?, ?, ? )`, t.Name(), t.OwnerID(), time.Now())
		if err != nil {
			slog.Error("create team repo", "error", err)
			return err
		}
		teamId, err = res.LastInsertId()
		if err != nil {
			slog.Error("last inserted id", "error", err)
			return err
		}
		teamMember, err := r.db.ExecContext(ctx, `
		INSERT INTO team_members (team_id, user_id, role, joined_at)
		VALUES (?, ?, ?, ? )`, teamId, t.OwnerID(), team.RoleOwner, time.Now())
		if err != nil {
			slog.Error("create team_member repo", "error", err)
			return err
		}
		idMem, err := teamMember.LastInsertId()
		if err != nil {
			slog.Error("last inserted id", "error", err)
			return err
		}
		slog.Info("team created", "id", teamId, "name", t.Name(), "owner_id", t.OwnerID())
		slog.Info("team_member created", "id", idMem, "team_id", teamId, "user_id", t.OwnerID(), "role", team.RoleOwner)

		return nil
	})
	if err != nil {
		return 0, err
	}
	return int64(teamId), err
}

func (r *TeamRepository) ListByUser(ctx context.Context, userID int64) ([]team.Team, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT t.id, t.name, t.owner_id, t.created_at
		FROM teams t
		JOIN team_members tm ON tm.team_id = t.id
		WHERE tm.user_id = ?
		ORDER BY t.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []team.Team
	for rows.Next() {
		var tr teamRow
		if err := rows.Scan(&tr.ID, &tr.Name, &tr.OwnerID, &tr.CreatedAt, &tr.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, toDomainTeam(tr))
	}
	return out, nil
}

func (r *TeamRepository) InviteUser(ctx context.Context, userID int64) error {
	return nil
}
