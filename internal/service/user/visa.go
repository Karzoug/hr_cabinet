package user

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) GetVisa(ctx context.Context, userID, passportID, visaID uint64) (*model.Visa, error) {
	const op = "user service: get visa"

	v, err := s.userRepository.GetVisa(ctx, userID, passportID, visaID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "visa not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return v, nil
}

func (s *service) ListVisas(ctx context.Context, userID, passportID uint64) ([]model.Visa, error) {
	const op = "user service: list visas"

	vs, err := s.userRepository.ListVisas(ctx, userID, passportID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return vs, nil
}

func (s *service) AddVisa(ctx context.Context, userID, passportID uint64, mv model.Visa) (uint64, error) {
	const op = "user service: add visa"

	id, err := s.userRepository.AddVisa(ctx, userID, passportID, mv)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: user/passport problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *service) UpdateVisa(ctx context.Context, userID, passportID uint64, v model.Visa) error {
	const op = "user service: update visa"

	err := s.userRepository.UpdateVisa(ctx, userID, passportID, v)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotAffected):
			return serr.NewError(serr.Conflict, "not updated: user/passport/visa problem")
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
