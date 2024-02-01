package user

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) Get(ctx context.Context, userID uint64) (*model.User, error) {
	const op = "user service: get user"

	u, err := s.userRepository.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "user not found")
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
			return nil, serr.NewError(serr.NotFound, "user not found")
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

	// TODO: add user to authorizations, use transaction

	id, err := s.userRepository.Add(ctx, u)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrConflict):
			return 0, serr.NewError(serr.Conflict, "not added: department or position not found")
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
		case errors.Is(err, repoerr.ErrRecordNotAffected):
			return serr.NewError(serr.Conflict, "not updated: user problem")
		case errors.Is(err, repoerr.ErrConflict):
			return serr.NewError(serr.Conflict, "not updated: department/position problem")
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
