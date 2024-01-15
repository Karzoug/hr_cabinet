package user

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) GetEducation(ctx context.Context, userID, educationID uint64) (*model.Education, error) {
	const op = "user service: get education"

	ed, err := s.userRepository.GetEducation(ctx, userID, educationID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "education not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ed, nil
}

func (s *service) ListEducations(ctx context.Context, userID uint64) ([]model.Education, error) {
	const op = "user service: list educations"

	eds, err := s.userRepository.ListEducations(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return eds, nil
}

func (s *service) AddEducation(ctx context.Context, userID uint64, ed model.Education) (uint64, error) {
	const op = "user service: add education"

	id, err := s.userRepository.AddEducation(ctx, userID, ed)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: user problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *service) UpdateEducation(ctx context.Context, userID uint64, ed model.Education) error {
	const op = "user service: update education"

	err := s.userRepository.UpdateEducation(ctx, userID, ed)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotAffected):
			return serr.NewError(serr.Conflict, "not updated: user/education problem")
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
