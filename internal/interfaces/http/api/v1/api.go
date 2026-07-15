package v1

type API struct {
	teamService        ITeamService
	authService        IUserService
	taskService        ITaskService
	taskHistoryService ITaskHistoryService
	commentService     ICommentService
}

func NewAPI(s ITeamService, u IUserService, t ITaskService, th ITaskHistoryService, c ICommentService) *API {
	return &API{
		teamService:        s,
		authService:        u,
		taskService:        t,
		taskHistoryService: th,
		commentService:     c,
	}
}
