package user

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

type passportRepo interface {
	List(ctx context.Context, userID uint64) ([]model.Passport, error)
	Get(ctx context.Context, userID, passportID uint64) (*model.Passport, error)
	Add(ctx context.Context, userID uint64, p model.Passport) (uint64, error)
	Update(ctx context.Context, userID uint64, p model.Passport) error
}

type PassportUseCase struct {
	repo passportRepo
}

func NewPassportUseCase(repo passportRepo) PassportUseCase {
	return PassportUseCase{
		repo: repo,
	}
}

func (s PassportUseCase) Get(ctx context.Context, userID, passportID uint64) (*model.Passport, error) {
	const op = "user service: get passport"

	p, err := s.repo.Get(ctx, userID, passportID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "passport not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return p, nil
}

func (s PassportUseCase) List(ctx context.Context, userID uint64) ([]model.Passport, error) {
	const op = "user service: list passports"

	psps, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return psps, nil
}

func (s PassportUseCase) Add(ctx context.Context, userID uint64, mp model.Passport) (uint64, error) {
	const op = "user service: add passport"

	id, err := s.repo.Add(ctx, userID, mp)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: department/position problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s PassportUseCase) Update(ctx context.Context, userID uint64, p model.Passport) error {
	const op = "user service: update passport"

	err := s.repo.Update(ctx, userID, p)
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
