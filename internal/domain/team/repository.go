package team

import "context"

type IRepository interface {
	Create(ctx context.Context, t *Team) (int64, error)
	ListByUser(ctx context.Context, userID int64) ([]Team, error)
	InviteUser(ctx context.Context, userID int64) error
}

type IMemberRepository interface {
	IsTeamMember(ctx context.Context, teamID, userID int64) (bool, error)
}
