package model

import (
	"time"
)

type User struct {
	ID                  uint64
	LastName            string
	FirstName           string
	MiddleName          string
	Gender              Gender
	DateOfBirth         time.Time
	PlaceOfBirth        string
	Grade               string
	PhoneNumbers        map[string]string
	Email               string
	RegistrationAddress string
	ResidentialAddress  string
	Nationality         string
	InsuranceNumber     string
	TaxpayerNumber      string
	Position            string
	Department          string
}

// Gender represents user gender.
type Gender string

const (
	GenderFemale Gender = "female"
	GenderMale   Gender = "male"
)
