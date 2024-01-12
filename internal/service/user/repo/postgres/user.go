package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/henvic/pgq"
	"github.com/jackc/pgx/v5"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

const LimitListUsers = 10

func (s *storage) Exist(ctx context.Context, userID uint64) (bool, error) {
	const op = "postrgresql user storage: exist user"

	row := s.DB.QueryRow(ctx, "SELECT COUNT(1) FROM users WHERE id = @user_id",
		pgx.NamedArgs{"user_id": userID})
	var count int
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if count == 0 {
		return false, nil
	}
	return true, nil
}

const getUserQuery = `SELECT 
users.id AS id,lastname,firstname,middlename,gender,
date_of_birth,place_of_birth,grade,phone_numbers,
work_email,registration_address,residential_address,nationality,
insurance_number,taxpayer_number, 
positions.title AS position, 
departments.title AS department 
FROM users		 
JOIN departments ON users.department_id = departments.id
JOIN positions ON users.position_id = positions.id
 WHERE users.id = @user_id`

func (s *storage) Get(ctx context.Context, userID uint64) (*model.User, error) {
	const op = "postrgresql user storage: get user"

	rows, err := s.DB.Query(ctx, getUserQuery, pgx.NamedArgs{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	u, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[user])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repoerr.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	mu := convertUserToModelUser(u)
	return &mu, nil
}

// GetExpandedUser returns summary information about the user.
// (!) This is a complex query that potentially returns a lot of data.
// Use context with timeout.
func (s *storage) GetExpandedUser(ctx context.Context, userID uint64) (*model.ExpandedUser, error) {
	const op = "postrgresql user storage: get expanded user"

	batch := &pgx.Batch{}
	batch.Queue(getUserQuery, pgx.NamedArgs{"user_id": userID})
	batch.Queue(listEducationsQuery, pgx.NamedArgs{"user_id": userID})
	batch.Queue(listTrainingsQuery, pgx.NamedArgs{"user_id": userID})
	batch.Queue(`SELECT 
	id, number, type, issued_date, issued_by	 
	FROM passports
	WHERE passports.user_id = @user_id`,
		pgx.NamedArgs{"user_id": userID})
	batch.Queue(`SELECT 
	id, number, passport_id, issued_state, 
	valid_to, valid_from, number_entries 
	FROM visas
	WHERE visas.user_id = @user_id`,
		pgx.NamedArgs{"user_id": userID})
	br := s.DB.SendBatch(ctx, batch)
	defer br.Close()

	var expUser model.ExpandedUser

	// TODO: duplicate code, refactor later maybe

	// get user
	rows, err := br.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	u, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[user])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repoerr.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	expUser.User = convertUserToModelUser(u)

	// get educations
	rows, err = br.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	eds, err := pgx.CollectRows[education](rows, pgx.RowToStructByNameLax[education])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	expUser.Educations = make([]model.Education, len(eds))
	for i, ed := range eds {
		expUser.Educations[i] = convertEducationToModelEducation(ed)
	}

	// get trainings
	rows, err = br.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	trs, err := pgx.CollectRows[training](rows, pgx.RowToStructByNameLax[training])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	expUser.Trainings = make([]model.Training, len(trs))
	for i, tr := range trs {
		expUser.Trainings[i] = convertTrainingToModelTraining(tr)
	}

	// get passports
	rows, err = br.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	psps, err := pgx.CollectRows[passport](rows, pgx.RowToStructByNameLax[passport])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	expUser.Passports = make([]model.PassportWithVisas, len(psps))
	for i, psp := range psps {
		expUser.Passports[i].Passport = convertPassportToModelPassport(psp)
	}

	// get visas
	rows, err = br.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	vs, err := pgx.CollectRows[visa](rows, pgx.RowToStructByNameLax[visa])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for i := 0; i < len(expUser.Passports); i++ {
		for j := 0; j < len(vs); j++ {
			if vs[j].PassportID == expUser.Passports[i].ID {
				expUser.Passports[i].Visas = append(expUser.Passports[i].Visas, convertVisaToModelVisa(vs[j]))
				expUser.Passports[i].VisasCount++
			}
		}
	}

	return &expUser, nil
}

func (s *storage) List(ctx context.Context, pms model.ListUsersParams) ([]model.User, int, error) {
	const op = "postrgresql user storage: list users"

	sb := pgq.
		Select(`users.id AS id,lastname,firstname,middlename,gender,
		date_of_birth,place_of_birth,grade,phone_numbers,
		work_email,registration_address,residential_address,nationality,
		insurance_number,taxpayer_number, 
		positions.title AS position, 
		departments.title AS department, count(*) OVER() AS total_count`).
		From("users").
		Join("departments ON users.department_id = departments.id").
		Join("positions ON users.position_id = positions.id")
	if pms.Query != "" {
		sb = sb.Where(pgq.ILike{"lastname": pms.Query + "%"})
	}
	// nolint:exhaustive
	switch pms.SortBy {
	case model.ListUsersParamsSortByDepartment:
		sb = sb.OrderBy("department, lastname, firstname")
	default:
		sb = sb.OrderBy("lastname, firstname")
	}
	sb = sb.
		Limit(uint64(pms.Limit)).
		Offset(uint64((pms.Page - 1) * pms.Limit))
	query, args, err := sb.SQL()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	lu, err := pgx.CollectRows[listUser](rows, pgx.RowToStructByNameLax[listUser])
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	if len(lu) == 0 {
		return []model.User{}, 0, nil
	}

	users := make([]model.User, len(lu))
	for i, u := range lu {
		users[i] = convertUserToModelUser(&u.user)
	}
	return users, lu[0].TotalCount, nil
}
