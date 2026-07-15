package dto

import "time"

type CommentListResult struct {
	Items []CommentItem `json:"items"`
	Page  int64         `json:"page"`
	Limit int64         `json:"limit"`
	Total int64         `json:"total"`
}

type CommentItem struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"task_id"`
	AuthorID  int64     `json:"author_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentListInput struct {
	TaskID int64
	UserID int64
	Page   int64
	Limit  int64
}
