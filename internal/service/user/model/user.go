package model

import (
	"time"
)

type ShortUserInfo struct {
	ID           uint64
	LastName     string
	FirstName    string
	MiddleName   string
	Department   string
	Position     string
	Email        string
	PhoneNumbers map[string]string
}

type Insurance struct {
	Number  string
	HasScan bool
}

type Taxpayer struct {
	Number  string
	HasScan bool
}

type Military struct {
	Rank         string
	Speciality   string
	Category     string
	Commissariat string
	HasScan      bool
}

type PersonalDataProcessing struct {
	HasScan bool
}

type User struct {
	ShortUserInfo
	Gender                 gender
	DateOfBirth            time.Time
	PlaceOfBirth           string
	Grade                  string
	RegistrationAddress    string
	ResidentialAddress     string
	Nationality            string
	Insurance              Insurance
	Taxpayer               Taxpayer
	PositionID             uint64
	DepartmentID           uint64
	Military               Military
	PersonalDataProcessing PersonalDataProcessing
}

// gender represents user gender.
type gender string

const (
	GenderFemale gender = "female"
	GenderMale   gender = "male"
)

// ExpandedUser represents summary information about the user.
type ExpandedUser struct {
	User
	Educations []Education
	Trainings  []Training
	Passports  []ExpandedPassport
	// Contracts []Contract
	Vacations []Vacation
}
