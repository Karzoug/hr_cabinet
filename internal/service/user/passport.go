package user

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) GetPassport(ctx context.Context, userID, passportID uint64) (*model.Passport, error) {
	const op = "user service: get passport"

	p, err := s.userRepository.GetPassport(ctx, userID, passportID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "passport not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return p, nil
}

func (s *service) ListPassports(ctx context.Context, userID uint64) ([]model.Passport, error) {
	const op = "user service: list passports"

	psps, err := s.userRepository.ListPassports(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return psps, nil
}

func (s *service) AddPassport(ctx context.Context, userID uint64, mp model.Passport) (uint64, error) {
	const op = "user service: add passport"

	id, err := s.userRepository.AddPassport(ctx, userID, mp)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: department/position problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *service) UpdatePassport(ctx context.Context, userID uint64, p model.Passport) error {
	const op = "user service: update passport"

	err := s.userRepository.UpdatePassport(ctx, userID, p)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotAffected):
			return serr.NewError(serr.Conflict, "not updated: user/passport problem")
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
