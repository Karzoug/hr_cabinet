package user

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

type contractRepo interface {
	List(ctx context.Context, userID uint64) ([]model.Contract, error)
	Get(ctx context.Context, userID, contractID uint64) (*model.Contract, error)
	Add(ctx context.Context, userID uint64, tr model.Contract) (uint64, error)
	Update(ctx context.Context, userID uint64, mc model.Contract) error
}

type ContractService struct {
	dbRepository contractRepo
}

func NewContractService(dbRepository contractRepo) ContractService {
	return ContractService{
		dbRepository: dbRepository,
	}
}

func (s ContractService) Get(ctx context.Context, userID, contractID uint64) (*model.Contract, error) {
	const op = "user service: get contract"

	tr, err := s.dbRepository.Get(ctx, userID, contractID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "contract not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tr, nil
}

func (s ContractService) List(ctx context.Context, userID uint64) ([]model.Contract, error) {
	const op = "user service: list contracts"

	ctrs, err := s.dbRepository.List(ctx, userID)
	if err != nil {
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	return ctrs, nil
}

func (s ContractService) Add(ctx context.Context, userID uint64, c model.Contract) (uint64, error) {
	const op = "user service: add contract"

	id, err := s.dbRepository.Add(ctx, userID, c)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return 0, serr.NewError(serr.Conflict, "not added: user problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s ContractService) Update(ctx context.Context, userID uint64, c model.Contract) error {
	const op = "user service: update contract"

	err := s.dbRepository.Update(ctx, userID, c)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotAffected):
			return serr.NewError(serr.Conflict, "not updated: user/contract problem")
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
