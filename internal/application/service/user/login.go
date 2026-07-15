package user

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Login(ctx context.Context, in LoginInput) (*LoginResult, error) {
	u, err := s.repo.ByEmail(ctx, in.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash()), []byte(in.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.token.Generate(u.ID(), u.Email())
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		UserID: u.ID(),
		Token:  token,
	}, nil
}
