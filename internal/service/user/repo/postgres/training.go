package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

const listTrainingsQuery = `SELECT 
id, title_of_program, title_of_institution, 
cost, date_end, date_begin
FROM trainings
WHERE user_id = @user_id`

func (s *storage) ListTrainings(ctx context.Context, userID uint64) ([]model.Training, error) {
	const op = "postrgresql user storage: list trainings"

	rows, err := s.DB.Query(ctx, listTrainingsQuery, pgx.NamedArgs{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	trs, err := pgx.CollectRows[training](rows, pgx.RowToStructByNameLax[training])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	trainings := make([]model.Training, len(trs))
	for i, tr := range trs {
		trainings[i] = convertTrainingToModelTraining(tr)
	}

	return trainings, nil
}

func (s *storage) GetTraining(ctx context.Context, userID, trainingID uint64) (*model.Training, error) {
	const op = "postrgresql user storage: get training"

	rows, err := s.DB.Query(ctx,
		`SELECT 
		id, title_of_program, title_of_institution, 
		cost, date_end, date_begin 
		FROM trainings
		WHERE id = @id AND user_id = @user_id`,
		pgx.NamedArgs{
			"id":      trainingID,
			"user_id": userID,
		})
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

	row := s.DB.QueryRow(ctx, `INSERT INTO trainings
		("user_id", "title_of_program", "title_of_institution", 
		"cost", "date_end", "date_begin")
		VALUES (@user_id, @title_of_program, @title_of_institution, 
		@cost, @date_end, @date_begin)
		RETURNING "id"`,
		pgx.NamedArgs{
			"user_id":              userID,
			"title_of_program":     tr.Program,
			"title_of_institution": tr.IssuedInstitution,
			"cost":                 tr.Cost,
			"date_end":             tr.DateTo,
			"date_begin":           tr.DateFrom,
		})

	if err := row.Scan(&tr.ID); err != nil {
		if strings.Contains(err.Error(), "23") && // Integrity Constraint Violation
			strings.Contains(err.Error(), "user_id") {
			return 0, fmt.Errorf("%s: the user does not exist: %w", op, repoerr.ErrRecordNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return tr.ID, nil
}

func (s *storage) UpdateTraining(ctx context.Context, userID uint64, tr model.Training) error {
	const op = "postrgresql user storage: update training"

	tag, err := s.DB.Exec(ctx, `UPDATE trainings
	SET title_of_program=@title_of_program, title_of_institution=@title_of_institution, 
	cost=@cost, date_end=@date_end, date_begin=@date_begin
	WHERE id=@id AND user_id=@user_id`,
		pgx.NamedArgs{
			"user_id":              userID,
			"id":                   tr.ID,
			"title_of_program":     tr.Program,
			"title_of_institution": tr.IssuedInstitution,
			"cost":                 tr.Cost,
			"date_end":             tr.DateTo,
			"date_begin":           tr.DateFrom,
		})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 { // it's ok for pgx
		return fmt.Errorf("%s: %w and %w", op,
			repoerr.ErrRecordNotModified,
			repoerr.ErrRecordNotFound)
	}
	return nil
}
