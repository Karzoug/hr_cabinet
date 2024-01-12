package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) GetVacation(ctx context.Context, userID, vacationID uint64) (*model.Vacation, error) {
	const op = "user service: get vacation"

	v, err := s.userRepository.GetVacation(ctx, userID, vacationID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrVacationNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return v, nil
}

func (s *service) ListVacations(ctx context.Context, userID uint64) ([]model.Vacation, error) {
	const op = "user service: list vacations"

	vcs, err := s.userRepository.ListVacations(ctx, userID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return vcs, nil
}

func (s *service) AddVacation(ctx context.Context, userID uint64, v model.Vacation) (uint64, error) {
	const op = "user service: add vacation"

	id, err := s.userRepository.AddVacation(ctx, userID, v)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return 0, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
