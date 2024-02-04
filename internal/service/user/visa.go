package user

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

type visaRepo interface {
	List(ctx context.Context, userID uint64) ([]model.Visa, error)
	Get(ctx context.Context, userID, visaID uint64) (*model.Visa, error)
	Add(ctx context.Context, userID uint64, mv model.Visa) (uint64, error)
	Update(ctx context.Context, userID uint64, v model.Visa) error
}

type VisaUseCase struct {
	visaRepo visaRepo
}

func NewVisaUseCase(visaRepo visaRepo) VisaUseCase {
	return VisaUseCase{
		visaRepo: visaRepo,
	}
}

func (s VisaUseCase) Get(ctx context.Context, userID, visaID uint64) (*model.Visa, error) {
	const op = "user VisaUseCase: get visa"

	v, err := s.visaRepo.Get(ctx, userID, visaID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "visa not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return v, nil
}

func (s VisaUseCase) List(ctx context.Context, userID uint64) ([]model.Visa, error) {
	const op = "user VisaUseCase: list visas"

	vs, err := s.visaRepo.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return vs, nil
}

func (s VisaUseCase) Add(ctx context.Context, userID uint64, mv model.Visa) (uint64, error) {
	const op = "user VisaUseCase: add visa"

	id, err := s.visaRepo.Add(ctx, userID, mv)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: user problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s VisaUseCase) Update(ctx context.Context, userID uint64, v model.Visa) error {
	const op = "user VisaUseCase: update visa"

	err := s.visaRepo.Update(ctx, userID, v)
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
