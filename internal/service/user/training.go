package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) GetTraining(ctx context.Context, trainingID uint64) (*model.Training, error) {
	const op = "user service: get training"

	ed, err := s.userRepository.GetTraining(ctx, trainingID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrTrainingNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ed, nil
}

func (s *service) ListTrainings(ctx context.Context, userID uint64) ([]model.Training, error) {
	const op = "user service: list trainings"

	eds, err := s.userRepository.ListTrainings(ctx, userID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return eds, nil
}

func (s *service) AddTraining(ctx context.Context, userID uint64, ed model.Training) (uint64, error) {
	const op = "user service: add training"

	id, err := s.userRepository.AddTraining(ctx, userID, ed)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return 0, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
