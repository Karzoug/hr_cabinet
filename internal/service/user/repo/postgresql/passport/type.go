package passport

import (
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

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

func convertFromDBO(p passport) model.Passport {
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

func convertToDBO(mp model.Passport) passport {
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
