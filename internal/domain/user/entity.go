package user

import "time"

type User struct {
	id           int64
	name         string
	email        string
	passwordHash string
	createdAt    time.Time
	updatedAt    time.Time
}

func New(name, email, passwordHash string) *User {
	now := time.Now()
	return &User{
		name:         name,
		email:        email,
		passwordHash: passwordHash,
		createdAt:    now,
		updatedAt:    now,
	}
}

func NewFromPersistence(id int64, name, email, passwordHash string, createdAt, updatedAt time.Time) *User {
	return &User{
		id:           id,
		name:         name,
		email:        email,
		passwordHash: passwordHash,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (u *User) SetID(id int64) {
	u.id = id
}

func (u *User) ID() int64            { return u.id }
func (u *User) Name() string         { return u.name }
func (u *User) Email() string        { return u.email }
func (u *User) PasswordHash() string { return u.passwordHash }
