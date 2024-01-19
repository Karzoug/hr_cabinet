package postgresql

import (
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type shortUserInfo struct {
	ID                uint64 `db:"id"`
	Department        string `db:"department"`
	Email             string `db:"work_email"`
	FirstName         string `db:"firstname"`
	LastName          string `db:"lastname"`
	MiddleName        string `db:"middlename"`
	MobilePhoneNumber string `db:"mobile_phone_number"`
	OfficePhoneNumber string `db:"office_phone_number"`
	Position          string `db:"position"`
}

type user struct {
	shortUserInfo
	Gender                        gender    `db:"gender"`
	DateOfBirth                   time.Time `db:"date_of_birth"`
	PlaceOfBirth                  string    `db:"place_of_birth"`
	Grade                         string    `db:"grade"`
	RegistrationAddress           string    `db:"registration_address"`
	ResidentialAddress            string    `db:"residential_address"`
	InsuranceNumber               string    `db:"insurance_number"`
	InsuranceHasScan              bool      `db:"insurance_has_scan"`
	TaxpayerNumber                string    `db:"taxpayer_number"`
	TaxpayerHasScan               bool      `db:"taxpayer_has_scan"`
	PositionID                    uint64    `db:"position_id"`
	DepartmentID                  uint64    `db:"department_id"`
	Military                      military
	PersonalDataProcessingHasScan bool `db:"pdp_has_scan"`
}

type gender string

const (
	genderFemale gender = "Женский"
	genderMale   gender = "Мужской"
)

type military struct {
	Rank         string `db:"rank"`
	Speciality   string `db:"specialty"`
	Category     string `db:"category_of_validity"`
	Commissariat string `db:"title_of_commissariat"`
	HasScan      bool   `db:"has_scan"`
}

func convertShortUserInfoToModelShortUserInfo(info shortUserInfo) model.ShortUserInfo {
	return model.ShortUserInfo(info)
}

func convertUserToModelUser(user *user) model.User {
	mu := model.User{
		ShortUserInfo:       convertShortUserInfoToModelShortUserInfo(user.shortUserInfo),
		DateOfBirth:         user.DateOfBirth,
		PlaceOfBirth:        user.PlaceOfBirth,
		Grade:               user.Grade,
		RegistrationAddress: user.RegistrationAddress,
		ResidentialAddress:  user.ResidentialAddress,
		Insurance: model.Insurance{
			Number:  user.InsuranceNumber,
			HasScan: user.InsuranceHasScan,
		},
		Taxpayer: model.Taxpayer{
			Number:  user.TaxpayerNumber,
			HasScan: user.TaxpayerHasScan,
		},
		PositionID:   user.PositionID,
		DepartmentID: user.DepartmentID,
		Military: model.Military{
			Rank:         user.Military.Rank,
			Speciality:   user.Military.Speciality,
			Category:     user.Military.Category,
			Commissariat: user.Military.Commissariat,
			HasScan:      user.Military.HasScan,
		},
		PersonalDataProcessing: model.PersonalDataProcessing{
			HasScan: user.PersonalDataProcessingHasScan,
		},
	}
	switch user.Gender {
	case genderMale:
		mu.Gender = model.GenderMale
	case genderFemale:
		mu.Gender = model.GenderFemale
	}
	return mu
}

func convertModelUserToUser(u *model.User) user {
	var gr gender
	switch u.Gender {
	case model.GenderMale:
		gr = genderMale
	case model.GenderFemale:
		gr = genderFemale
	}

	return user{
		shortUserInfo:       shortUserInfo(u.ShortUserInfo),
		Gender:              gr,
		DateOfBirth:         u.DateOfBirth,
		PlaceOfBirth:        u.PlaceOfBirth,
		Grade:               u.Grade,
		RegistrationAddress: u.RegistrationAddress,
		ResidentialAddress:  u.ResidentialAddress,
		InsuranceNumber:     u.Insurance.Number,
		TaxpayerNumber:      u.Taxpayer.Number,
		PositionID:          u.PositionID,
		DepartmentID:        u.DepartmentID,
	}
}

type listUser struct {
	shortUserInfo
	TotalCount int `db:"total_count"`
}

type education struct {
	ID                uint64    `db:"id"`
	Number            string    `db:"document_number"`
	Program           string    `db:"title_of_program"`
	IssuedInstitution string    `db:"title_of_institution"`
	DateTo            time.Time `db:"year_of_end"`
	DateFrom          time.Time `db:"year_of_begin"`
	HasScan           bool      `db:"has_scan"`
}

func convertEducationToModelEducation(ed education) model.Education {
	return model.Education(ed)
}

type training struct {
	ID                uint64    `db:"id"`
	Program           string    `db:"title_of_program"`
	IssuedInstitution string    `db:"title_of_institution"`
	Cost              uint64    `db:"cost"`
	DateTo            time.Time `db:"date_end"`
	DateFrom          time.Time `db:"date_begin"`
	HasScan           bool      `db:"has_scan"`
}

func convertTrainingToModelTraining(tr training) model.Training {
	return model.Training(tr)
}

type passport struct {
	ID           uint64       `db:"id"`
	Citizenship  string       `db:"citizenship"`
	IssuedBy     *string      `db:"issued_by"`
	IssuedByCode *string      `db:"issued_by_code"`
	IssuedDate   time.Time    `db:"issued_date"`
	Number       string       `db:"number"`
	Type         passportType `db:"type"`
	HasScan      bool         `db:"has_scan"`
}

type passportType string

const (
	passportTypeInternational passportType = "Заграничный"
	passportTypeNational      passportType = "Внутренний"
)

func convertPassportToModelPassport(p passport) model.Passport {
	var pt model.PassportType
	switch p.Type {
	case passportTypeInternational:
		pt = model.PassportTypeInternational
	case passportTypeNational:
		pt = model.PassportTypeNational
	}

	return model.Passport{
		ID:          p.ID,
		Citizenship: p.Citizenship,
		IssuedBy:    p.IssuedBy,
		IssuedDate:  p.IssuedDate,
		Number:      p.Number,
		Type:        pt,
		HasScan:     p.HasScan,
	}
}

func convertModelPassportToPassport(mp model.Passport) passport {
	var t passportType
	switch mp.Type {
	case model.PassportTypeInternational:
		t = passportTypeInternational
	case model.PassportTypeNational:
		t = passportTypeNational
	}

	return passport{
		ID:          mp.ID,
		Citizenship: mp.Citizenship,
		IssuedBy:    mp.IssuedBy,
		IssuedDate:  mp.IssuedDate,
		Number:      mp.Number,
		Type:        t,
	}
}

type visa struct {
	ID          uint64    `db:"id"`
	Number      string    `db:"number"`
	Type        string    `db:"type"`
	IssuedState *string   `db:"issued_state"`
	ValidTo     time.Time `db:"valid_to"`
	ValidFrom   time.Time `db:"valid_from"`
}

func convertVisaToModelVisa(v visa) model.Visa {
	return model.Visa{
		ID:          v.ID,
		Number:      v.Number,
		Type:        model.VisaType(v.Type),
		IssuedState: v.IssuedState,
		ValidTo:     v.ValidTo,
		ValidFrom:   v.ValidFrom,
	}
}

func convertModelVisaToVisa(mv model.Visa) visa {
	return visa{
		ID:          mv.ID,
		Number:      mv.Number,
		Type:        string(mv.Type),
		IssuedState: mv.IssuedState,
		ValidTo:     mv.ValidTo,
		ValidFrom:   mv.ValidFrom,
	}
}

type vacation struct {
	ID        uint64    `db:"id"`
	DateBegin time.Time `db:"date_begin"`
	DateEnd   time.Time `db:"date_end"`
}

func convertVacationToModelVacation(v vacation) model.Vacation {
	return model.Vacation(v)
}

type scan struct {
	ID          uint64    `db:"id"`
	UserID      uint64    `db:"user_id"`
	Type        scanType  `db:"type"`
	DocumentID  uint64    `db:"document_id"`
	Description string    `db:"description"`
	UploadedAt  time.Time `db:"created_at"`
}

type scanType string

const (
	scanTypePassport   scanType = "Паспорт"
	scanTypeTaxpayer   scanType = "ИНН"
	scanTypeInsurance  scanType = "СНИЛС"
	scanTypeContract   scanType = "Трудовой договор"
	scanTypePDP        scanType = "Согласие на обработку данных"
	scanTypeMilitary   scanType = "Военный билет"
	scanTypeEducation  scanType = "Документ об образовании"
	scanTypeTraining   scanType = "Сертификат"
	scanTypeBriefing   scanType = "Инструктаж"
	scanTypeWorkPermit scanType = "Разрешение на работу"
	scanTypeMarriage   scanType = "Свидетельство о браке"
	scanTypeBabyBirth  scanType = "Свидетельство о рождении"
	scanTypeOther      scanType = "Другое"
)

func convertScanToModelScan(s scan) model.Scan {
	var st model.ScanType
	switch s.Type {
	case scanTypePassport:
		st = model.ScanTypePassport
	case scanTypeTaxpayer:
		st = model.ScanTypeTaxpayer
	case scanTypeInsurance:
		st = model.ScanTypeInsurance
	case scanTypeContract:
		st = model.ScanTypeContract
	case scanTypePDP:
		st = model.ScanTypePDP
	case scanTypeMilitary:
		st = model.ScanTypeMilitary
	case scanTypeEducation:
		st = model.ScanTypeEducation
	case scanTypeTraining:
		st = model.ScanTypeTraining
	case scanTypeBriefing:
		st = model.ScanTypeBriefing
	case scanTypeWorkPermit:
		st = model.ScanTypeWorkPermit
	case scanTypeMarriage:
		st = model.ScanTypeMarriage
	case scanTypeBabyBirth:
		st = model.ScanTypeBabyBirth
	case scanTypeOther:
		st = model.ScanTypeOther
	}

	return model.Scan{
		ID:          s.ID,
		Type:        st,
		DocumentID:  s.DocumentID,
		Description: s.Description,
		UploadedAt:  s.UploadedAt,
	}
}

func convertModelScanToScan(ms model.Scan) scan {
	var t scanType
	switch ms.Type {
	case model.ScanTypePassport:
		t = scanTypePassport
	case model.ScanTypeTaxpayer:
		t = scanTypeTaxpayer
	case model.ScanTypeInsurance:
		t = scanTypeInsurance
	case model.ScanTypeContract:
		t = scanTypeContract
	case model.ScanTypePDP:
		t = scanTypePDP
	case model.ScanTypeMilitary:
		t = scanTypeMilitary
	case model.ScanTypeEducation:
		t = scanTypeEducation
	case model.ScanTypeTraining:
		t = scanTypeTraining
	case model.ScanTypeBriefing:
		t = scanTypeBriefing
	case model.ScanTypeWorkPermit:
		t = scanTypeWorkPermit
	case model.ScanTypeMarriage:
		t = scanTypeMarriage
	case model.ScanTypeBabyBirth:
		t = scanTypeBabyBirth
	case model.ScanTypeOther:
		t = scanTypeOther
	}

	return scan{
		ID:          ms.ID,
		DocumentID:  ms.DocumentID,
		Type:        t,
		Description: ms.Description,
	}
}

type contract struct {
	ID              uint64       `db:"id"`
	Number          string       `db:"number"`
	ContractType    contractType `db:"contract_type"`
	WorkTypeID      uint64       `db:"work_type_id"`
	ProbationPeriod *uint        `db:"probation_period"`
	DateBegin       time.Time    `db:"date_begin"`
	DateEnd         *time.Time   `db:"date_end"`
	HasScan         bool         `db:"has_scan"`
}

type contractType string

const (
	contractTypePermanent contractType = "Бессрочный"
	contractTypeTemporary contractType = "Срочный"
)

func convertContractToModelContract(c contract) model.Contract {
	mc := model.Contract{
		ID:              c.ID,
		Number:          c.Number,
		WorkTypeID:      c.WorkTypeID,
		ProbationPeriod: c.ProbationPeriod,
		DateBegin:       c.DateBegin,
		DateEnd:         c.DateEnd,
		HasScan:         c.HasScan,
	}

	switch c.ContractType {
	case contractTypePermanent:
		mc.Type = model.ContractTypePermanent
	case contractTypeTemporary:
		mc.Type = model.ContractTypeTemporary
	}

	return mc
}

func convertModelContractToContract(mc model.Contract) contract {
	c := contract{
		ID:              mc.ID,
		Number:          mc.Number,
		WorkTypeID:      mc.WorkTypeID,
		ProbationPeriod: mc.ProbationPeriod,
		DateBegin:       mc.DateBegin,
		DateEnd:         mc.DateEnd,
	}

	switch mc.Type {
	case model.ContractTypePermanent:
		c.ContractType = contractTypePermanent
	case model.ContractTypeTemporary:
		c.ContractType = contractTypeTemporary
	}

	return c
}
