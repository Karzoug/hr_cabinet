package visa

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

type service struct {
	dbRepository dbRepository
}

func NewService(dbRepository dbRepository) *service {
	return &service{
		dbRepository: dbRepository,
	}
}

func (s *service) Get(ctx context.Context, userID, visaID uint64) (*model.Visa, error) {
	const op = "visa service: get"

	v, err := s.dbRepository.Get(ctx, userID, visaID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "visa not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return v, nil
}

func (s *service) List(ctx context.Context, userID uint64) ([]model.Visa, error) {
	const op = "visa service: list"

	vs, err := s.dbRepository.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return vs, nil
}

func (s *service) Add(ctx context.Context, userID uint64, mv model.Visa) (uint64, error) {
	const op = "visa service: add"

	id, err := s.dbRepository.Add(ctx, userID, mv)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: user problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *service) Update(ctx context.Context, userID uint64, v model.Visa) error {
	const op = "visa service: update"

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
