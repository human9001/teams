package input

import (
	"time"

	"github.com/human9001/teams/internal/domain/task/model"
)

type ListTasksInput struct {
	TeamID     int64         `json:"team_id"`
	Status     *model.Status `json:"status"`
	AssigneeID *int64        `json:"asignee_id"`
	Page       int64         `json:"page"`
	Limit      int64         `json:"limit"`
}

type ListTasksResult struct {
	Items []ListTasksItem `json:"items"`
	Total int64           `json:"total"`
	Page  int64           `json:"page"`
	Limit int64           `json:"limit"`
}

type ListTasksItem struct {
	ID          int64     `json:"id"`
	TeamID      int64     `json:"team_id"`
	CreatedBy   int64     `json:"user_id"`
	AssigneeID  int64     `json:"asignee_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateTaskInput struct {
	TaskID      int64
	UserID      int64
	Title       *string
	Description *string
	Status      *model.Status
	Priority    *model.Priority
	AssigneeID  *int64
}

type CreateTaskInput struct {
	TeamID      int64
	AssigneeID  *int64
	Title       string
	Description string
	Priority    model.Priority
}

type ListFilter struct {
	TeamID     int64
	Status     *model.Status
	AssigneeID *int64
	Page       int64
	Limit      int64
}

type ListFilterResult struct {
	Items []ListTasksItem
	Total int64
}
