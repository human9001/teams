package model

import "time"

type EventType string

const (
	EventCreated EventType = "created"
	EventUpdated EventType = "updated"
	EventDeleted EventType = "deleted"
)

type HistoryItem struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"task_id"`
	EventType EventType `json:"event_type"`
	ActorID   *int64    `json:"actor_id,omitempty"`
	Before    []byte    `json:"before,omitempty"`
	After     []byte    `json:"after,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
