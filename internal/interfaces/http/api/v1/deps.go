package v1

import (
	"context"

	commentDto "github.com/human9001/teams/internal/application/service/comment/dto"
	"github.com/human9001/teams/internal/application/service/history/dto"
	"github.com/human9001/teams/internal/application/service/task/input"
	"github.com/human9001/teams/internal/application/service/user"
	"github.com/human9001/teams/internal/domain/task/model"
	"github.com/human9001/teams/internal/domain/team"
	domainUser "github.com/human9001/teams/internal/domain/user"
)

type ITeamService interface {
	CreateTeam(ctx context.Context, name string, ownerID int64) (*team.Team, error)
	ListTeams(ctx context.Context, ownerID int64) ([]team.Team, error)
}

type IUserService interface {
	Login(ctx context.Context, in user.LoginInput) (*user.LoginResult, error)
	Register(ctx context.Context, in user.RegisterInput) (*domainUser.User, error)
}

type ITaskService interface {
	CreateTask(ctx context.Context, userId int64, req input.CreateTaskInput) (*model.Task, error)
	ListTasks(ctx context.Context, in input.ListTasksInput) (input.ListTasksResult, error)
	UpdateTask(ctx context.Context, in input.UpdateTaskInput) error
}

type ITaskHistoryService interface {
	GetTaskHistory(ctx context.Context, taskID int64) (dto.ListHistoryResult, error)
}

type ICommentService interface {
	ListComments(ctx context.Context, in commentDto.CommentListInput) (commentDto.CommentListResult, error)
}
