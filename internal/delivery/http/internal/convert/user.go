package convert

import (
	"github.com/oapi-codegen/runtime/types"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

func FromAPIAddUserRequest(req api.AddUserJSONRequestBody) model.User {
	user := model.User{
		ShortUserInfo: model.ShortUserInfo{
			Email:      string(req.Email),
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			MiddleName: req.MiddleName,
		},
		DateOfBirth:         req.DateOfBirth.Time,
		PlaceOfBirth:        req.PlaceOfBirth,
		Grade:               req.Grade,
		RegistrationAddress: req.RegistrationAddress,
		ResidentialAddress:  req.ResidentialAddress,
		Nationality:         req.Nationality,
		InsuranceNumber:     req.Insurance.Number,
		TaxpayerNumber:      req.Taxpayer.Number,
		PositionID:          req.PositionID,
		DepartmentID:        req.DepartmentID,
	}
	switch req.Gender {
	case api.Female:
		user.Gender = model.GenderFemale
	case api.Male:
		user.Gender = model.GenderMale
	}
	if req.PhoneNumbers != nil {
		user.PhoneNumbers = make(map[string]string, len(req.PhoneNumbers))
		for k, v := range req.PhoneNumbers {
			user.PhoneNumbers[k] = string(v)
		}
	}
	return user
}

func FromAPIPutUserRequest(userID uint64, req api.PutUserJSONRequestBody) model.User {
	user := model.User{
		ShortUserInfo: model.ShortUserInfo{
			ID:         userID,
			Email:      string(req.Email),
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			MiddleName: req.MiddleName,
		},
		DateOfBirth:         req.DateOfBirth.Time,
		PlaceOfBirth:        req.PlaceOfBirth,
		Grade:               req.Grade,
		RegistrationAddress: req.RegistrationAddress,
		ResidentialAddress:  req.ResidentialAddress,
		Nationality:         req.Nationality,
		InsuranceNumber:     req.Insurance.Number,
		TaxpayerNumber:      req.Taxpayer.Number,
		PositionID:          req.PositionID,
		DepartmentID:        req.DepartmentID,
	}
	switch req.Gender {
	case api.Female:
		user.Gender = model.GenderFemale
	case api.Male:
		user.Gender = model.GenderMale
	}
	if req.PhoneNumbers != nil {
		user.PhoneNumbers = make(map[string]string, len(req.PhoneNumbers))
		for k, v := range req.PhoneNumbers {
			user.PhoneNumbers[k] = string(v)
		}
	}
	return user
}

func ToAPIGetExpandedUserResponse(u *model.ExpandedUser) api.GetExpandedUserResponse {
	var expUser api.GetExpandedUserResponse

	expUser.GetUserResponse = ToAPIGetUserResponse(&u.User)
	expUser.Educations = ToAPIListEducations(u.Educations)
	expUser.Trainings = ToAPIListTrainings(u.Trainings)
	expUser.Passports = ToAPIExpandedPassports(u.Passports)
	expUser.Vacations = ToAPIListVacations(u.Vacations)

	// TODO: add contracts and vacation
	expUser.Contracts = []api.Contract{}

	return expUser
}

func ToAPIGetUserResponse(u *model.User) api.GetUserResponse {
	resp := api.GetUserResponse{
		DateOfBirth:         types.Date{Time: u.DateOfBirth},
		DepartmentID:        u.DepartmentID,
		Email:               types.Email(u.Email),
		FirstName:           u.FirstName,
		Gender:              "",
		Grade:               u.Grade,
		ID:                  u.ID,
		Insurance:           api.Insurance{Number: u.InsuranceNumber},
		LastName:            u.LastName,
		MiddleName:          u.MiddleName,
		Nationality:         u.Nationality,
		PlaceOfBirth:        u.PlaceOfBirth,
		PositionID:          u.PositionID,
		RegistrationAddress: u.RegistrationAddress,
		ResidentialAddress:  u.ResidentialAddress,
		Taxpayer:            api.Taxpayer{Number: u.TaxpayerNumber},
	}
	if u.PhoneNumbers != nil {
		resp.PhoneNumbers = make(map[string]api.PhoneNumber, len(u.PhoneNumbers))
		for k, v := range u.PhoneNumbers {
			resp.PhoneNumbers[k] = api.PhoneNumber(v)
		}
	}
	switch u.Gender {
	case model.GenderFemale:
		resp.Gender = api.Female
	case model.GenderMale:
		resp.Gender = api.Male
	}
	return resp
}

func ToAPIListUsers(users []model.ShortUserInfo) []api.ListUsersItem {
	res := make([]api.ListUsersItem, len(users))
	for i := 0; i < len(users); i++ {
		res[i] = toAPIListUser(users[i])
	}
	return res
}

func toAPIListUser(u model.ShortUserInfo) api.ListUsersItem {
	item := api.ListUsersItem{
		ID:         u.ID,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		MiddleName: u.MiddleName,
		Email:      types.Email(u.Email),
		Position:   u.Position,
		Department: u.Department,
	}
	if u.PhoneNumbers != nil {
		item.PhoneNumbers = make(map[string]api.PhoneNumber, len(u.PhoneNumbers))
		for k, v := range u.PhoneNumbers {
			item.PhoneNumbers[k] = api.PhoneNumber(v)
		}
	}
	return item
}
