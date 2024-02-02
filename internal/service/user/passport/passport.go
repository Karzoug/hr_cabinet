package passport

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

func (s *service) Get(ctx context.Context, userID, passportID uint64) (*model.Passport, error) {
	const op = "passport service: get"

	p, err := s.dbRepository.Get(ctx, userID, passportID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "passport not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return p, nil
}

func (s *service) List(ctx context.Context, userID uint64) ([]model.Passport, error) {
	const op = "passport service: list"

	psps, err := s.dbRepository.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return psps, nil
}

func (s *service) Add(ctx context.Context, userID uint64, mp model.Passport) (uint64, error) {
	const op = "passport service: add"

	id, err := s.dbRepository.Add(ctx, userID, mp)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: department/position problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *service) Update(ctx context.Context, userID uint64, p model.Passport) error {
	const op = "passport service: update"

	err := s.dbRepository.Update(ctx, userID, p)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotAffected):
			return serr.NewError(serr.Conflict, "not updated: user/passport problem")
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
