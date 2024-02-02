package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/henvic/pgq"
	"github.com/jackc/pgx/v5"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	pq "github.com/Employee-s-file-cabinet/backend/pkg/postgresql"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

const LimitListUsers = 10

type storage struct {
	pq.DB
}

func New(db pq.DB) storage {
	return storage{DB: db}
}

func (s storage) Exist(ctx context.Context, userID uint64) (bool, error) {
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

const (
	getUserQuery = `SELECT 
users.id AS id, lastname, firstname, middlename, gender,
date_of_birth, place_of_birth, grade, mobile_phone_number, office_phone_number,
work_email, registration_address, residential_address,
insurance_number, taxpayer_number, users.department_id AS department_id, position_id,
positions.title AS position, departments.title AS department,
(SELECT COUNT(*)>0 FROM scans WHERE user_id=@user_id AND scans.type='ИНН') AS insurance_has_scan,
(SELECT COUNT(*)>0 FROM scans WHERE user_id=@user_id AND scans.type='СНИЛС') AS taxpayer_has_scan,
(SELECT COUNT(*)>0 FROM scans WHERE user_id=@user_id AND scans.type='Согласие на обработку данных') AS pdp_has_scan
FROM users
JOIN departments ON users.department_id = departments.id
JOIN positions ON users.position_id = positions.id
WHERE users.id = @user_id`

	getMilitaryQuery = `SELECT
rank, specialty, category_of_validity, title_of_commissariat,
(SELECT COUNT(*)>0 FROM scans WHERE scans.user_id=@user_id AND scans.document_id=militaries.id AND scans.type='Военный билет') AS has_scan
FROM militaries
WHERE militaries.user_id = @user_id`
)

func (s storage) Get(ctx context.Context, userID uint64) (*model.User, error) {
	const op = "postgresql user storage: get user"

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

	rows, err = s.DB.Query(ctx, getMilitaryQuery, pgx.NamedArgs{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	m, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[military])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repoerr.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	u.Military = *m
	mu := convertFromDBO(u)
	return &mu, nil
}

// // GetExpandedUser returns summary information about the user.
// // (!) This is a complex query that potentially returns a lot of data.
// // Use context with timeout.
// func (s storage) GetExpandedUser(ctx context.Context, userID uint64) (*model.ExpandedUser, error) {
// 	const op = "postgresql user storage: get expanded user"

// 	batch := &pgx.Batch{}
// 	batch.Queue(getUserQuery, pgx.NamedArgs{"user_id": userID})
// 	batch.Queue(listEducationsQuery, pgx.NamedArgs{"user_id": userID})
// 	batch.Queue(listTrainingsQuery, pgx.NamedArgs{"user_id": userID})
// 	batch.Queue(listPassportsQuery, pgx.NamedArgs{"user_id": userID})
// 	batch.Queue(listVisasQuery, pgx.NamedArgs{"user_id": userID})
// 	batch.Queue(listVacationsQuery, pgx.NamedArgs{"user_id": userID})
// 	batch.Queue(listContractsQuery, pgx.NamedArgs{"user_id": userID})
// 	br := s.DB.SendBatch(ctx, batch)
// 	defer br.Close()

// 	var expUser model.ExpandedUser

// 	// TODO: duplicate code, refactor later maybe

// 	// get user
// 	rows, err := br.Query()
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	u, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[user])
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return nil, repoerr.ErrRecordNotFound
// 		}
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	expUser.User = convertUserToModelUser(u)

// 	// get educations
// 	rows, err = br.Query()
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	eds, err := pgx.CollectRows[education](rows, pgx.RowToStructByNameLax[education])
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	expUser.Educations = make([]model.Education, len(eds))
// 	for i, ed := range eds {
// 		expUser.Educations[i] = convertEducationToModelEducation(ed)
// 	}

// 	// get trainings
// 	rows, err = br.Query()
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	trs, err := pgx.CollectRows[training](rows, pgx.RowToStructByNameLax[training])
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	expUser.Trainings = make([]model.Training, len(trs))
// 	for i, tr := range trs {
// 		expUser.Trainings[i] = convertTrainingToModelTraining(tr)
// 	}

// 	// get passports
// 	rows, err = br.Query()
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	psps, err := pgx.CollectRows[passport](rows, pgx.RowToStructByNameLax[passport])
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	expUser.Passports = make([]model.Passport, len(psps))
// 	for i, psp := range psps {
// 		expUser.Passports[i] = convertPassportToModelPassport(psp)
// 	}

// 	// get visas
// 	rows, err = br.Query()
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	visas, err := pgx.CollectRows[visa](rows, pgx.RowToStructByNameLax[visa])
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	expUser.Visas = make([]model.Visa, len(visas))
// 	for i, v := range visas {
// 		expUser.Visas[i] = convertVisaToModelVisa(v)
// 	}

// 	// get vacations
// 	rows, err = br.Query()
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	vs, err := pgx.CollectRows[vacation](rows, pgx.RowToStructByNameLax[vacation])
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	expUser.Vacations = make([]model.Vacation, len(vs))
// 	for i, v := range vs {
// 		expUser.Vacations[i] = convertVacationToModelVacation(v)
// 	}

// 	// get contracts
// 	rows, err = br.Query()
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	cs, err := pgx.CollectRows[contract](rows, pgx.RowToStructByNameLax[contract])
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	expUser.Contracts = make([]model.Contract, len(cs))
// 	for i, c := range cs {
// 		expUser.Contracts[i] = convertContractToModelContract(c)
// 	}

// 	return &expUser, nil
// }

func (s storage) ListShortUserInfo(ctx context.Context, pms model.ListUsersParams) ([]model.ShortUserInfo, int, error) {
	const op = "postgresql user storage: list short user info"

	sb := pgq.
		Select(`users.id AS id, lastname, firstname, middlename, 
		mobile_phone_number, office_phone_number, work_email, 
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
		users[i] = convertFromShortUserInfoDBO(u.shortUserInfo)
	}
	return users, lu[0].TotalCount, nil
}

func (s storage) Add(ctx context.Context, mu model.User) (uint64, error) {
	const op = "postrgresql user storage: add user"

	user := convertToDBO(&mu)

	row := s.DB.QueryRow(ctx,
		`INSERT INTO users 
			(lastname, firstname, middlename, 
			gender, date_of_birth, place_of_birth, 
			grade, mobile_phone_number, office_phone_number, work_email, 
			registration_address, residential_address, 
			insurance_number, taxpayer_number, 
			department_id, position_id)
		VALUES
			(@lastname, @firstname, @middlename, 
			@gender, @date_of_birth, @place_of_birth, 
			@grade, @mobile_phone_number, @office_phone_number, @email, 
			@registration_address, @residential_address, 
			@insurance_number, @taxpayer_number, 
			@department_id, @position_id)
			RETURNING id`,
		pgx.NamedArgs{
			"lastname":             user.LastName,
			"firstname":            user.FirstName,
			"middlename":           user.MiddleName,
			"gender":               user.Gender,
			"date_of_birth":        user.DateOfBirth,
			"place_of_birth":       user.PlaceOfBirth,
			"grade":                user.Grade,
			"mobile_phone_number":  user.MobilePhoneNumber,
			"office_phone_number":  user.OfficePhoneNumber,
			"email":                user.Email,
			"registration_address": user.RegistrationAddress,
			"residential_address":  user.ResidentialAddress,
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

func (s storage) Update(ctx context.Context, mu model.User) error {
	const op = "postrgresql user storage: update user"

	user := convertToDBO(&mu)

	tag, err := s.DB.Exec(ctx, `UPDATE users
	SET lastname = @lastname, firstname = @firstname, middlename = @middlename, 
	gender = @gender, date_of_birth = @date_of_birth, place_of_birth = @place_of_birth, 
	grade = @grade, mobile_phone_number = @mobile_phone_number, office_phone_number = @office_phone_number, work_email = @email, 
	registration_address = @registration_address, residential_address = @residential_address, 
	insurance_number = @insurance_number, taxpayer_number = @taxpayer_number, 
	department_id = @department_id, position_id = @position_id
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
			"mobile_phone_number":  user.MobilePhoneNumber,
			"office_phone_number":  user.OfficePhoneNumber,
			"email":                user.Email,
			"registration_address": user.RegistrationAddress,
			"residential_address":  user.ResidentialAddress,
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
