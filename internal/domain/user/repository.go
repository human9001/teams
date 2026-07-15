package user

import "context"

type IRepository interface {
	Create(ctx context.Context, u *User) (int64, error)
	ByEmail(ctx context.Context, email string) (*User, error)
}
