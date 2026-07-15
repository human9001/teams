package model

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"task_id"`
	AuthorID  int64     `json:"author_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type ListResult struct {
	Items []Comment `json:"items"`
	Page  int64     `json:"page"`
	Limit int64     `json:"limit"`
	Total int64     `json:"total"`
}

type ListInput struct {
	TaskID int64
	UserID int64
	Page   int64
	Limit  int64
}
