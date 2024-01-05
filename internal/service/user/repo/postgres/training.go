package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/henvic/pgq"
	"github.com/jackc/pgx/v5"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *storage) ListTrainings(ctx context.Context, userID uint64) ([]model.Training, error) {
	const op = "postrgresql user storage: list trainings"

	rows, err := s.DB.Query(ctx, `SELECT 
	id, title_of_program, title_of_institution, 
	cost, date_end, date_begin
	FROM trainings
	WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	eds, err := pgx.CollectRows[training](rows, pgx.RowToStructByNameLax[training])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	trainings := make([]model.Training, len(eds))
	for i, ed := range eds {
		trainings[i] = convertTrainingToModelTraining(ed)
	}

	return trainings, nil
}

func (s *storage) GetTraining(ctx context.Context, trainingID uint64) (*model.Training, error) {
	const op = "postrgresql user storage: get training"

	rows, err := s.DB.Query(ctx,
		`SELECT 
		id, title_of_program, title_of_institution, 
		cost, date_end, date_begin 
		FROM trainings
		WHERE id = $1`, trainingID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	ed, err := pgx.CollectExactlyOneRow[training](rows, pgx.RowToStructByNameLax[training])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, repoerr.ErrRecordNotFound)
	}

	med := convertTrainingToModelTraining(ed)
	return &med, nil
}

func (s *storage) AddTraining(ctx context.Context, userID uint64, tr model.Training) (uint64, error) {
	const op = "postrgresql user storage: add training"

	qb := pgq.Insert("trainings").
		Columns("user_id", "title_of_program", "title_of_institution", "cost", "date_end", "date_begin").
		Values(userID, tr.Program, tr.IssuedInstitution, tr.Cost, tr.DateTo, tr.DateFrom).
		Returning("id")
	query, args, err := qb.SQL()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.DB.QueryRow(ctx, query, args...).Scan(&tr.ID); err != nil {
		if strings.Contains(err.Error(), "23") && // Integrity Constraint Violation
			strings.Contains(err.Error(), "user_id") {
			return 0, fmt.Errorf("%s: the user does not exist: %w", op, repoerr.ErrRecordNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return tr.ID, nil
}
