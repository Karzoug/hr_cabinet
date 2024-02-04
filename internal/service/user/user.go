package user

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

type userRepo interface {
	Exist(ctx context.Context, userID uint64) (bool, error)
	ListShortUserInfo(ctx context.Context, pms model.ListUsersParams) ([]model.ShortUserInfo, int, error)
	Get(ctx context.Context, userID uint64) (*model.User, error)
	GetExpandedUser(ctx context.Context, userID uint64) (*model.ExpandedUser, error)
	Add(ctx context.Context, user model.User) (uint64, error)
	Update(ctx context.Context, user model.User) error
}

type UserUseCase struct {
	userRepo userRepo
}

func NewUserUseCase(userRepo userRepo) UserUseCase {
	return UserUseCase{
		userRepo: userRepo,
	}
}

func (s UserUseCase) Get(ctx context.Context, userID uint64) (*model.User, error) {
	const op = "user service: get user"

	u, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "user not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return u, nil
}

func (s UserUseCase) GetExpanded(ctx context.Context, userID uint64) (*model.ExpandedUser, error) {
	const op = "user service: get expanded"

	eu, err := s.userRepo.GetExpandedUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "user not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return eu, nil
}

func (s UserUseCase) ListShortUserInfo(ctx context.Context, params model.ListUsersParams) ([]model.ShortUserInfo, int, error) {
	const op = "user service: list users"

	users, count, err := s.userRepo.ListShortUserInfo(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	return users, count, nil
}

func (s UserUseCase) Add(ctx context.Context, u model.User) (uint64, error) {
	const op = "user service: add user"

	// TODO: add user to authorizations, use transaction

	id, err := s.userRepo.Add(ctx, u)
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

func (s UserUseCase) Update(ctx context.Context, user model.User) error {
	const op = "user service: update user"

	err := s.userRepo.Update(ctx, user)
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
