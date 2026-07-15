package team

import (
	"github.com/human9001/teams/internal/domain/team"
)

type TeamService struct {
	repo team.IRepository
}

func NewTeamService(repo team.IRepository) *TeamService {
	return &TeamService{
		repo: repo,
	}
}
