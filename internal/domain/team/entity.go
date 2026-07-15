package team

import (
	"errors"
	"strings"
	"time"
)

type Role string

const (
	RoleOwner  Role = "OWNER"
	RoleAdmin  Role = "ADMIN"
	RoleMember Role = "MEMBER"
)

func (r Role) Valid() bool {
	switch r {
	case RoleOwner, RoleAdmin, RoleMember:
		return true
	default:
		return false
	}
}

type Team struct {
	id        int64
	name      string
	ownerID   int64
	createdAt time.Time
}

func New(name string, ownerID int64) (*Team, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("team name is required")
	}
	if ownerID == 0 {
		return nil, errors.New("owner id is required")
	}

	now := time.Now()
	return &Team{
		name:      name,
		ownerID:   ownerID,
		createdAt: now,
	}, nil
}

func (t *Team) SetID(id int64) {
	t.id = id
}

func (t *Team) Invite(userID int64) error {
	if userID == 0 {
		return errors.New("user id is required")
	}
	return nil
}

func NewFromPersistence(id int64, name string, ownerID int64, createdAt, updatedAt time.Time) Team {
	return Team{id: id, name: name, ownerID: ownerID, createdAt: createdAt}
}

func (t *Team) ID() int64            { return t.id }
func (t *Team) Name() string         { return t.name }
func (t *Team) OwnerID() int64       { return t.ownerID }
func (t *Team) CreatedAt() time.Time { return t.createdAt }
