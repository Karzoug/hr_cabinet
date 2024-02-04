package user

import (
	"context"
	"errors"
	"fmt"

	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

type educationRepo interface {
	Get(ctx context.Context, userID, educationID uint64) (*model.Education, error)
	List(ctx context.Context, userID uint64) ([]model.Education, error)
	Add(ctx context.Context, userID uint64, ed model.Education) (uint64, error)
	Update(ctx context.Context, userID uint64, ed model.Education) error
}

type EducationUseCase struct {
	repo educationRepo
}

func NewEducationUseCase(repo educationRepo) EducationUseCase {
	return EducationUseCase{
		repo: repo,
	}
}

func (s EducationUseCase) Get(ctx context.Context, userID, educationID uint64) (*model.Education, error) {
	const op = "user service: get education"

	ed, err := s.repo.Get(ctx, userID, educationID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "education not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ed, nil
}

func (s EducationUseCase) List(ctx context.Context, userID uint64) ([]model.Education, error) {
	const op = "user service: list educations"

	eds, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return eds, nil
}

func (s EducationUseCase) Add(ctx context.Context, userID uint64, ed model.Education) (uint64, error) {
	const op = "user service: add education"

	id, err := s.repo.Add(ctx, userID, ed)
	if err != nil {
		if errors.Is(err, repoerr.ErrConflict) {
			return 0, serr.NewError(serr.Conflict, "not added: user problem")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s EducationUseCase) Update(ctx context.Context, userID uint64, ed model.Education) error {
	const op = "user service: update education"

	err := s.repo.Update(ctx, userID, ed)
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
