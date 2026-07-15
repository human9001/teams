package user

import (
	"github.com/human9001/teams/internal/domain/user"
)

type AuthService struct {
	repo  user.IRepository
	token TokenGenerator
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	UserID int64
	Token  string
}

type TokenGenerator interface {
	Generate(userID int64, email string) (string, error)
}

func NewAuthService(r user.IRepository, t TokenGenerator) *AuthService {
	return &AuthService{repo: r, token: t}
}
