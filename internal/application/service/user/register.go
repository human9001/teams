package user

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/human9001/teams/internal/domain/user"
)

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*user.User, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	if in.Name == "" || in.Email == "" || in.Password == "" {
		return nil, errors.New("name, email and password are required")
	}

	existing, err := s.repo.ByEmail(ctx, in.Email)
	if err != nil && err.Error() != "not found" {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := user.New(in.Name, in.Email, string(hash))
	userId, err := s.repo.Create(ctx, u)
	if err != nil {
		return nil, err
	}
	u.SetID(userId)
	return u, nil
}
