package user

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

type trainingRepo interface {
	List(ctx context.Context, userID uint64) ([]model.Training, error)
	Get(ctx context.Context, userID, trainingID uint64) (*model.Training, error)
	Add(ctx context.Context, userID uint64, tr model.Training) (uint64, error)
	Update(ctx context.Context, userID uint64, tr model.Training) error
}

type TrainingUseCase struct {
	repo trainingRepo
}

func NewTrainingUseCase(repo trainingRepo) TrainingUseCase {
	return TrainingUseCase{
		repo: repo,
	}
}

func (s TrainingUseCase) Get(ctx context.Context, userID, trainingID uint64) (*model.Training, error) {
	const op = "user service: get training"

	tr, err := s.repo.Get(ctx, userID, trainingID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "training not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tr, nil
}

func (s TrainingUseCase) List(ctx context.Context, userID uint64) ([]model.Training, error) {
	const op = "user service: list trainings"

	trs, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return trs, nil
}

func (s TrainingUseCase) Add(ctx context.Context, userID uint64, ed model.Training) (uint64, error) {
	const op = "user service: add training"

	id, err := s.repo.Add(ctx, userID, ed)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: user problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s TrainingUseCase) Update(ctx context.Context, userID uint64, tr model.Training) error {
	const op = "user service: update training"

	err := s.repo.Update(ctx, userID, tr)
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
