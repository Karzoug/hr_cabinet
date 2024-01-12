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

func (s *service) GetExpanded(ctx context.Context, userID uint64) (*model.ExpandedUser, error) {
	const op = "user service: get expanded"

	eu, err := s.userRepository.GetExpandedUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return eu, nil
}

func (s *service) ListShortUserInfo(ctx context.Context, params model.ListUsersParams) ([]model.ShortUserInfo, int, error) {
	const op = "user service: list users"

	users, count, err := s.userRepository.ListShortUserInfo(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	return users, count, nil
}

func (s *service) Add(ctx context.Context, u model.User) (uint64, error) {
	const op = "user service: add user"

	id, err := s.userRepository.Add(ctx, u)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotFound):
			return 0, fmt.Errorf("%s: %w", op, ErrDepartmentOrPositionNotFound)
		default:
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}
	return id, nil
}

func (s *service) Update(ctx context.Context, user model.User) error {
	const op = "user service: update user"

	err := s.userRepository.Update(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotModified):
			return ErrUserNotFound
		case errors.Is(err, repoerr.ErrRecordNotFound):
			return ErrDepartmentOrPositionNotFound
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
