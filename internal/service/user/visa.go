package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) GetVisa(ctx context.Context, userID, passportID, visaID uint64) (*model.Visa, error) {
	const op = "user service: get visa"

	v, err := s.userRepository.GetVisa(ctx, userID, passportID, visaID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrVisaNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return v, nil
}

func (s *service) ListVisas(ctx context.Context, userID, passportID uint64) ([]model.Visa, error) {
	const op = "user service: list visas"

	vs, err := s.userRepository.ListVisas(ctx, userID, passportID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserOrPassportNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return vs, nil
}

func (s *service) AddVisa(ctx context.Context, userID, passportID uint64, mv model.Visa) (uint64, error) {
	const op = "user service: add visa"

	id, err := s.userRepository.AddVisa(ctx, userID, passportID, mv)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return 0, fmt.Errorf("%s: %w", op, ErrUserOrPassportNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
