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
