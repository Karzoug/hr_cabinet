package convert

import (
	"github.com/oapi-codegen/runtime/types"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

func ToAPIExpandedFullUser(u *model.ExpandedUser) api.ExpandedFullUser {
	var expUser api.ExpandedFullUser

	expUser.FullUser = ToAPIFullUser(&u.User)
	expUser.Educations = ToAPIEducations(u.Educations)
	expUser.Trainings = ToAPITrainings(u.Trainings)
	expUser.Passports = ToAPIPassportsWithVisas(u.Passports)

	// TODO: add contracts and vacation
	expUser.Contracts = []api.Contract{}
	expUser.Vacations = []api.Vacation{}

	return expUser
}

func ToAPIFullUser(u *model.User) api.FullUser {
	var gr api.Gender
	switch u.Gender {
	case model.GenderFemale:
		gr = api.GenderFemale
	case model.GenderMale:
		gr = api.GenderMale
	}
	return api.FullUser{
		ShortUser: ToAPIShortUser(u),
		DateOfBirth: types.Date{
			Time: u.DateOfBirth,
		},
		Gender: gr,
		Grade:  u.Grade,
		Insurance: api.Insurance{
			Number: u.InsuranceNumber,
		},
		Nationality:         u.Nationality,
		PlaceOfBirth:        u.PlaceOfBirth,
		RegistrationAddress: u.RegistrationAddress,
		ResidentialAddress:  u.ResidentialAddress,
		Taxpayer: api.Taxpayer{
			Number: u.TaxpayerNumber,
		},
	}
}

func ToAPIShortUsers(users []model.User) []api.ShortUser {
	res := make([]api.ShortUser, len(users))
	for i := 0; i < len(users); i++ {
		res[i] = ToAPIShortUser(&users[i])
	}
	return res
}

func ToAPIShortUser(u *model.User) api.ShortUser {
	var phn map[string]api.PhoneNumber
	if u.PhoneNumbers != nil {
		phn = make(map[string]api.PhoneNumber, len(u.PhoneNumbers))
		for k, v := range u.PhoneNumbers {
			phn[k] = api.PhoneNumber(v)
		}
	}
	return api.ShortUser{
		ID:           &u.ID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		MiddleName:   u.MiddleName,
		Email:        u.Email,
		Position:     u.Position,
		Department:   u.Department,
		PhoneNumbers: phn,
	}
}
