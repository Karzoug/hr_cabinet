package visa

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

func (s Subservice) Get(ctx context.Context, userID, visaID uint64) (*model.Visa, error) {
	const op = "user Subservice: get visa"

	v, err := s.dbRepository.Get(ctx, userID, visaID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "visa not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return v, nil
}

func (s Subservice) List(ctx context.Context, userID uint64) ([]model.Visa, error) {
	const op = "user Subservice: list visas"

	vs, err := s.dbRepository.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return vs, nil
}

func (s Subservice) Add(ctx context.Context, userID uint64, mv model.Visa) (uint64, error) {
	const op = "user Subservice: add visa"

	id, err := s.dbRepository.Add(ctx, userID, mv)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: user problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s Subservice) Update(ctx context.Context, userID uint64, v model.Visa) error {
	const op = "user Subservice: update visa"

	err := s.dbRepository.Update(ctx, userID, v)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotAffected):
			return serr.NewError(serr.Conflict, "not updated: user/visa problem")
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
