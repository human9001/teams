package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type MembershipRepository struct {
	db *sqlx.DB
}

func NewMembershipRepository(db *sqlx.DB) *MembershipRepository {
	return &MembershipRepository{db: db}
}

func (r *MembershipRepository) IsTeamMember(ctx context.Context, teamID, userID int64) (bool, error) {
	var x int
	err := r.db.GetContext(ctx, &x, `
		SELECT 1 FROM team_members
		WHERE team_id = ? AND user_id = ?
		LIMIT 1
	`, teamID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
