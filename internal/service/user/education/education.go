package education

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

func (s Subservice) Get(ctx context.Context, userID, educationID uint64) (*model.Education, error) {
	const op = "user service: get education"

	ed, err := s.dbRepository.Get(ctx, userID, educationID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "education not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ed, nil
}

func (s Subservice) List(ctx context.Context, userID uint64) ([]model.Education, error) {
	const op = "user service: list educations"

	eds, err := s.dbRepository.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return eds, nil
}

func (s Subservice) Add(ctx context.Context, userID uint64, ed model.Education) (uint64, error) {
	const op = "user service: add education"

	id, err := s.dbRepository.Add(ctx, userID, ed)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: user problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s Subservice) Update(ctx context.Context, userID uint64, ed model.Education) error {
	const op = "user service: update education"

	err := s.dbRepository.Update(ctx, userID, ed)
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
