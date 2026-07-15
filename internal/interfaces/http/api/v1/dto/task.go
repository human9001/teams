package dto

type CreateTaskRequest struct {
	TeamID      int64  `json:"team_id"`
	AssigneeID  *int64 `json:"assignee_id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

type TaskResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

type ListTasksResponse struct {
	Items []TaskItem `json:"items"`
	Page  int64      `json:"page"`
	Limit int64      `json:"limit"`
	Total int64      `json:"total"`
}

type TaskItem struct {
	ID          int64  `json:"id"`
	TeamID      int64  `json:"team_id"`
	AssigneeID  *int64 `json:"assignee_id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	Priority    *string `json:"priority"`
	AssigneeID  *int64  `json:"assignee_id"`
	DueDate     *string `json:"due_date"`
}
