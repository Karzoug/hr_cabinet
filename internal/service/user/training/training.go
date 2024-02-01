package training

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) Get(ctx context.Context, userID, trainingID uint64) (*model.Training, error) {
	const op = "training service: get"

	tr, err := s.dbRepository.Get(ctx, userID, trainingID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "training not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tr, nil
}

func (s *service) List(ctx context.Context, userID uint64) ([]model.Training, error) {
	const op = "training service: list"

	trs, err := s.dbRepository.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return trs, nil
}

func (s *service) Add(ctx context.Context, userID uint64, ed model.Training) (uint64, error) {
	const op = "training service: add"

	id, err := s.dbRepository.Add(ctx, userID, ed)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: user problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *service) Update(ctx context.Context, userID uint64, tr model.Training) error {
	const op = "training service: update"

	err := s.dbRepository.Update(ctx, userID, tr)
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
