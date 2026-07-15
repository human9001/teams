package model

import (
	"errors"
	"time"
)

type Task struct {
	id          int64
	teamID      int64
	createdBy   int64
	assigneeID  *int64
	title       string
	description string
	status      Status
	priority    Priority
	createdAt   time.Time
	updatedAt   time.Time
}

func New(teamID, createdBy int64, title, description string, priority Priority) (*Task, error) {
	if teamID == 0 {
		return nil, errors.New("team id is required")
	}
	if createdBy == 0 {
		return nil, errors.New("created by is required")
	}
	if title == "" {
		return nil, errors.New("title is required")
	}
	if !priority.Valid() {
		return nil, errors.New("invalid priority")
	}

	now := time.Now()
	return &Task{
		teamID:      teamID,
		createdBy:   createdBy,
		title:       title,
		description: description,
		status:      StatusNew,
		priority:    priority,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func (t *Task) Assign(userID int64) error {
	if userID == 0 {
		return errors.New("invalid user id")
	}
	t.assigneeID = &userID
	t.updatedAt = time.Now()
	return nil
}

func (t *Task) ChangeStatus(newStatus Status) error {
	if !newStatus.Valid() {
		return errors.New("invalid status")
	}
	if !t.canTransitionTo(newStatus) {
		return errors.New("invalid status transition")
	}
	t.status = newStatus
	t.updatedAt = time.Now()
	return nil
}

func (t *Task) UpdateDetails(title, description string, priority Priority) error {
	if title == "" {
		return errors.New("title is required")
	}
	if !priority.Valid() {
		return errors.New("invalid priority")
	}
	t.title = title
	t.description = description
	t.priority = priority
	t.updatedAt = time.Now()
	return nil
}

func (t *Task) canTransitionTo(newStatus Status) bool {
	switch t.status {
	case StatusNew:
		return newStatus == StatusInProgress || newStatus == StatusCancelled
	case StatusInProgress:
		return newStatus == StatusDone || newStatus == StatusCancelled
	case StatusDone, StatusCancelled:
		return false
	default:
		return false
	}
}

func (t *Task) SetID(id int64) {
	t.id = id
}

func (t *Task) ID() int64            { return t.id }
func (t *Task) TeamID() int64        { return t.teamID }
func (t *Task) CreatedBy() int64     { return t.createdBy }
func (t *Task) AssigneeID() *int64   { return t.assigneeID }
func (t *Task) Title() string        { return t.title }
func (t *Task) Description() string  { return t.description }
func (t *Task) Status() Status       { return t.status }
func (t *Task) Priority() Priority   { return t.priority }
func (t *Task) CreatedAt() time.Time { return t.createdAt }
func (t *Task) UpdatedAt() time.Time { return t.updatedAt }

func (t *Task) SetTitle(v string)       { t.title = v }
func (t *Task) SetDescription(v string) { t.description = v }
func (t *Task) SetStatus(v Status)      { t.status = v }
func (t *Task) SetPriority(v Priority)  { t.priority = v }
func (t *Task) SetAssigneeID(v *int64)  { t.assigneeID = v }

func NewFromPersistence(id, teamID, createdBy int64, assigneeID *int64, title, description string,
	status Status, priority Priority, createdAt, updatedAt time.Time,
) Task {
	return Task{
		id:          id,
		teamID:      teamID,
		createdBy:   createdBy,
		assigneeID:  assigneeID,
		title:       title,
		description: description,
		status:      status,
		priority:    priority,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
