package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

const (
	MaxPhotoSize  = 20 << 20 // bytes
	photoFileName = "photo"
)

func (s *service) Get(ctx context.Context, userID uint64) (*model.User, error) {
	const op = "user service: get user"

	u, err := s.userRepository.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return u, nil
}

func (s *service) List(ctx context.Context, params model.ListUsersParams) ([]model.User, int, error) {
	const op = "user service: list users"

	users, count, err := s.userRepository.List(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	return users, count, nil
}
