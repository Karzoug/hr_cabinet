package vacation

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

func (s Subservice) GetVacation(ctx context.Context, userID, vacationID uint64) (*model.Vacation, error) {
	const op = "user service: get vacation"

	v, err := s.dbRepository.Get(ctx, userID, vacationID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "vacation not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return v, nil
}

func (s Subservice) ListVacations(ctx context.Context, userID uint64) ([]model.Vacation, error) {
	const op = "user service: list vacations"

	vcs, err := s.dbRepository.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return vcs, nil
}

func (s Subservice) Add(ctx context.Context, userID uint64, v model.Vacation) (uint64, error) {
	const op = "user service: add vacation"

	id, err := s.dbRepository.Add(ctx, userID, v)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: user problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s Subservice) Update(ctx context.Context, userID uint64, v model.Vacation) error {
	const op = "user service: update vacation"

	err := s.dbRepository.Update(ctx, userID, v)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotAffected):
			return serr.NewError(serr.Conflict, "not updated: vacation/user problem")
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
