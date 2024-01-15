package user

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) GetTraining(ctx context.Context, userID, trainingID uint64) (*model.Training, error) {
	const op = "user service: get training"

	tr, err := s.userRepository.GetTraining(ctx, userID, trainingID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "training not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tr, nil
}

func (s *service) ListTrainings(ctx context.Context, userID uint64) ([]model.Training, error) {
	const op = "user service: list trainings"

	trs, err := s.userRepository.ListTrainings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return trs, nil
}

func (s *service) AddTraining(ctx context.Context, userID uint64, ed model.Training) (uint64, error) {
	const op = "user service: add training"

	id, err := s.userRepository.AddTraining(ctx, userID, ed)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: user problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *service) UpdateTraining(ctx context.Context, userID uint64, tr model.Training) error {
	const op = "user service: update training"

	err := s.userRepository.UpdateTraining(ctx, userID, tr)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotAffected):
			return serr.NewError(serr.Conflict, "not updated: user/training problem")
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
