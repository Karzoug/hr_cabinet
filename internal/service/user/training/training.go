package training

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

type Subservice struct {
	dbRepository dbRepository
}

func New(dbRepository dbRepository) Subservice {
	return Subservice{
		dbRepository: dbRepository,
	}
}

func (s Subservice) Get(ctx context.Context, userID, trainingID uint64) (*model.Training, error) {
	const op = "user service: get training"

	tr, err := s.dbRepository.Get(ctx, userID, trainingID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "training not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tr, nil
}

func (s Subservice) List(ctx context.Context, userID uint64) ([]model.Training, error) {
	const op = "user service: list trainings"

	trs, err := s.dbRepository.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return trs, nil
}

func (s Subservice) Add(ctx context.Context, userID uint64, ed model.Training) (uint64, error) {
	const op = "user service: add training"

	id, err := s.dbRepository.Add(ctx, userID, ed)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: user problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s Subservice) Update(ctx context.Context, userID uint64, tr model.Training) error {
	const op = "user service: update training"

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
