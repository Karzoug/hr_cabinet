package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/service/auth/model/token"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) Login(ctx context.Context, login, password string) (string, error) {
	const op = "auth service: login"

	authnData, err := s.authRepository.Get(ctx, login)
	if errors.Is(err, repoerr.ErrRecordNotFound) {
		return "", fmt.Errorf("%s: %w", op, ErrForbidden)
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.passwordVerificator.Check(password, authnData.PasswordHash)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrForbidden)
	}

	t, err := s.tokenManager.Create(
		token.Data{
			UserID: authnData.UserID,
			RoleID: authnData.RoleID,
		})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}

func (s *service) Expires() time.Time {
	return s.tokenManager.Expires()
}
