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
users.id AS id, firstname, middlename, lastname, gender,
date_of_birth, place_of_birth, grade, phone_numbers,
work_email, registration_address,residential_address, nationality,
insurance_number, taxpayer_number, users.department_id AS department_id, position_id,
departments.title AS department, positions.title AS position
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
			return nil, repoerr.ErrRecordNotFound
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
	batch.Queue(listVacationsQuery, pgx.NamedArgs{"user_id": userID})
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
			return nil, repoerr.ErrRecordNotFound
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
	expUser.Passports = make([]model.ExpandedPassport, len(psps))
	for i, psp := range psps {
		expUser.Passports[i].Passport = convertPassportToModelPassport(psp)
	}

	// get visas
	rows, err = br.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	visas, err := pgx.CollectRows[visa](rows, pgx.RowToStructByNameLax[visa])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for i := 0; i < len(expUser.Passports); i++ {
		for j := 0; j < len(visas); j++ {
			if visas[j].PassportID == expUser.Passports[i].ID {
				expUser.Passports[i].Visas = append(expUser.Passports[i].Visas, convertVisaToModelVisa(visas[j]))
				expUser.Passports[i].VisasCount++
			}
		}
	}

	// get vacations
	rows, err = br.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	vs, err := pgx.CollectRows[vacation](rows, pgx.RowToStructByNameLax[vacation])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	expUser.Vacations = make([]model.Vacation, len(vs))
	for i, v := range vs {
		expUser.Vacations[i] = convertVacationToModelVacation(v)
	}

	return &expUser, nil
}

func (s *storage) ListShortUserInfo(ctx context.Context, pms model.ListUsersParams) ([]model.ShortUserInfo, int, error) {
	const op = "postrgresql user storage: list short user info"

	sb := pgq.
		Select(`users.id AS id, lastname, firstname, middlename, 
		phone_numbers, work_email, 
		positions.title AS position, departments.title AS department, 
		count(*) OVER() AS total_count`).
		From("users").
		Join("departments ON users.department_id = departments.id").
		Join("positions ON users.position_id = positions.id")
	if pms.Query != "" {
		q := "%" + pms.Query + "%"
		sb = sb.Where(`users.lastname ILIKE ? OR departments.title ILIKE ?`, q, q)
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
		return []model.ShortUserInfo{}, 0, nil
	}

	users := make([]model.ShortUserInfo, len(lu))
	for i, u := range lu {
		users[i] = convertShortUserInfoToModelShortUserInfo(u.shortUserInfo)
	}
	return users, lu[0].TotalCount, nil
}

func (s *storage) Add(ctx context.Context, mu model.User) (uint64, error) {
	const op = "postrgresql user storage: add user"

	user := convertModelUserToUser(&mu)

	row := s.DB.QueryRow(ctx,
		`INSERT INTO users 
			(lastname, firstname, middlename, 
			gender, date_of_birth, place_of_birth, 
			grade, phone_numbers, work_email, 
			registration_address, residential_address, 
			nationality, insurance_number, 
			taxpayer_number, department_id, position_id)
		VALUES
			(@lastname, @firstname, @middlename, 
			@gender, @date_of_birth, @place_of_birth, 
			@grade, @phone_numbers, @email, 
			@registration_address, @residential_address, 
			@nationality, @insurance_number, 
			@taxpayer_number, @department_id, @position_id)
			RETURNING id`,
		pgx.NamedArgs{
			"lastname":             user.LastName,
			"firstname":            user.FirstName,
			"middlename":           user.MiddleName,
			"gender":               user.Gender,
			"date_of_birth":        user.DateOfBirth,
			"place_of_birth":       user.PlaceOfBirth,
			"grade":                user.Grade,
			"phone_numbers":        user.PhoneNumbers,
			"email":                user.Email,
			"registration_address": user.RegistrationAddress,
			"residential_address":  user.ResidentialAddress,
			"nationality":          user.Nationality,
			"insurance_number":     user.InsuranceNumber,
			"taxpayer_number":      user.TaxpayerNumber,
			"department_id":        user.DepartmentID,
			"position_id":          user.PositionID,
		})

	if err := row.Scan(&user.ID); err != nil {
		if strings.Contains(err.Error(), "23") { // Integrity Constraint Violation
			if strings.Contains(err.Error(), "department_id") {
				return 0, fmt.Errorf("the department does not exist: %w", repoerr.ErrConflict)
			}
			if strings.Contains(err.Error(), "passport_id") {
				return 0, fmt.Errorf("the position does not exist: %w", repoerr.ErrConflict)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return user.ID, nil
}

func (s *storage) Update(ctx context.Context, mu model.User) error {
	const op = "postrgresql user storage: update user"

	user := convertModelUserToUser(&mu)

	tag, err := s.DB.Exec(ctx, `UPDATE users
	SET lastname=@lastname, firstname=@firstname, middlename=@middlename, 
	gender=@gender, date_of_birth=@date_of_birth, place_of_birth=@place_of_birth, 
	grade=@grade, phone_numbers=@phone_numbers, work_email=@email, 
	registration_address=@registration_address, residential_address=@residential_address, 
	nationality=@nationality, insurance_number=@insurance_number, 
	taxpayer_number=@taxpayer_number, 
	department_id=@department_id, position_id=@position_id
	WHERE id=@id`,
		pgx.NamedArgs{
			"id":                   user.ID,
			"lastname":             user.LastName,
			"firstname":            user.FirstName,
			"middlename":           user.MiddleName,
			"gender":               user.Gender,
			"date_of_birth":        user.DateOfBirth,
			"place_of_birth":       user.PlaceOfBirth,
			"grade":                user.Grade,
			"phone_numbers":        user.PhoneNumbers,
			"email":                user.Email,
			"registration_address": user.RegistrationAddress,
			"residential_address":  user.ResidentialAddress,
			"nationality":          user.Nationality,
			"insurance_number":     user.InsuranceNumber,
			"taxpayer_number":      user.TaxpayerNumber,
			"department_id":        user.DepartmentID,
			"position_id":          user.PositionID,
		})

	if err != nil {
		if strings.Contains(err.Error(), "23") { // Integrity Constraint Violation
			if strings.Contains(err.Error(), "department_id") {
				return fmt.Errorf("the department does not exist: %w", repoerr.ErrConflict)
			}
			if strings.Contains(err.Error(), "passport_id") {
				return fmt.Errorf("the position does not exist: %w", repoerr.ErrConflict)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 { // it's ok for pgx
		return repoerr.ErrRecordNotAffected
	}
	return nil
}
