package user

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) GetContract(ctx context.Context, userID, contractID uint64) (*model.Contract, error) {
	const op = "user service: get contract"

	tr, err := s.userRepository.GetContract(ctx, userID, contractID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "contract not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tr, nil
}

func (s *service) ListContracts(ctx context.Context, userID uint64) ([]model.Contract, error) {
	const op = "user service: list contracts"

	ctrs, err := s.userRepository.ListContracts(ctx, userID)
	if err != nil {
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	return ctrs, nil
}

func (s *service) AddContract(ctx context.Context, userID uint64, c model.Contract) (uint64, error) {
	const op = "user service: add contract"

	id, err := s.userRepository.AddContract(ctx, userID, c)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return 0, serr.NewError(serr.Conflict, "not added: user problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *service) UpdateContract(ctx context.Context, userID uint64, c model.Contract) error {
	const op = "user service: update contract"

	err := s.userRepository.UpdateContract(ctx, userID, c)
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
