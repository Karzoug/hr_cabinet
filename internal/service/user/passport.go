package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) GetPassport(ctx context.Context, userID, passportID uint64) (*model.Passport, error) {
	const op = "user service: get passport"

	p, err := s.userRepository.GetPassport(ctx, userID, passportID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrPassportNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return p, nil
}

func (s *service) ListPassports(ctx context.Context, userID uint64) ([]model.Passport, error) {
	const op = "user service: list passports"

	psps, err := s.userRepository.ListPassports(ctx, userID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return psps, nil
}

func (s *service) AddPassport(ctx context.Context, userID uint64, mp model.Passport) (uint64, error) {
	const op = "user service: add passport"

	id, err := s.userRepository.AddPassport(ctx, userID, mp)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return 0, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
