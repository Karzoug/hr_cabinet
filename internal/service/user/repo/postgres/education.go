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

func (s *storage) ListEducations(ctx context.Context, userID uint64) ([]model.Education, error) {
	const op = "postrgresql user storage: list educations"

	rows, err := s.DB.Query(ctx, `SELECT 
	id, document_number, title_of_program, 
	title_of_institution, year_of_end, year_of_begin 
	FROM educations
	WHERE user_id = @user_id`,
		pgx.NamedArgs{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	eds, err := pgx.CollectRows[education](rows, pgx.RowToStructByNameLax[education])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	educations := make([]model.Education, len(eds))
	for i, ed := range eds {
		educations[i] = convertEducationToModelEducation(ed)
	}

	return educations, nil
}

func (s *storage) GetEducation(ctx context.Context, userID, educationID uint64) (*model.Education, error) {
	const op = "postrgresql user storage: get education"

	rows, err := s.DB.Query(ctx, `SELECT 
		id, document_number, title_of_program, 
		title_of_institution, year_of_end, year_of_begin 
		FROM educations
		WHERE id = @education_id AND user_id = @user_id`,
		pgx.NamedArgs{
			"education_id": educationID,
			"user_id":      userID,
		})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	ed, err := pgx.CollectExactlyOneRow[education](rows, pgx.RowToStructByNameLax[education])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, repoerr.ErrRecordNotFound)
	}

	med := convertEducationToModelEducation(ed)
	return &med, nil
}

func (s *storage) AddEducation(ctx context.Context, userID uint64, ed model.Education) (uint64, error) {
	const op = "postrgresql user storage: add education"

	row := s.DB.QueryRow(ctx, `INSERT INTO educations
		("user_id", "document_number", "title_of_program", 
		"title_of_institution", "year_of_end", "year_of_begin") 
		VALUES (@user_id, @education_number, @education_program, 
		@education_issued_institution, @education_date_to, @education_date_from)
		RETURNING "id"`,
		pgx.NamedArgs{
			"user_id":                      userID,
			"education_number":             ed.Number,
			"education_program":            ed.Program,
			"education_issued_institution": ed.IssuedInstitution,
			"education_date_to":            ed.DateTo,
			"education_date_from":          ed.DateFrom,
		})

	if err := row.Scan(&ed.ID); err != nil {
		if strings.Contains(err.Error(), "23") && // Integrity Constraint Violation
			strings.Contains(err.Error(), "user_id") {
			return 0, fmt.Errorf("%s: the user does not exist: %w", op, repoerr.ErrRecordNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return ed.ID, nil
}
