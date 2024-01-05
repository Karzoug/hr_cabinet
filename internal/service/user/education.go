package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) GetEducation(ctx context.Context, educationID uint64) (*model.Education, error) {
	const op = "user service: get education"

	ed, err := s.userRepository.GetEducation(ctx, educationID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrEducationNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ed, nil
}

func (s *service) ListEducations(ctx context.Context, userID uint64) ([]model.Education, error) {
	const op = "user service: list educations"

	eds, err := s.userRepository.ListEducations(ctx, userID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return eds, nil
}

func (s *service) AddEducation(ctx context.Context, userID uint64, ed model.Education) (uint64, error) {
	const op = "user service: add education"

	id, err := s.userRepository.AddEducation(ctx, userID, ed)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return 0, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
