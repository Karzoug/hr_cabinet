package user

import (
	"github.com/oapi-codegen/runtime/types"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

func fromAPIAddUserRequest(req api.AddUserJSONRequestBody) model.User {
	user := model.User{
		ShortUserInfo: model.ShortUserInfo{
			Email:             string(req.Email),
			FirstName:         req.FirstName,
			LastName:          req.LastName,
			MiddleName:        req.MiddleName,
			MobilePhoneNumber: req.MobilePhoneNumber,
			OfficePhoneNumber: req.OfficePhoneNumber,
		},
		DateOfBirth:         req.DateOfBirth.Time,
		PlaceOfBirth:        req.PlaceOfBirth,
		Grade:               req.Grade,
		RegistrationAddress: req.RegistrationAddress,
		ResidentialAddress:  req.ResidentialAddress,
		Insurance:           model.Insurance{Number: req.Insurance.Number},
		Taxpayer:            model.Taxpayer{Number: req.Taxpayer.Number},
		PositionID:          req.PositionID,
		DepartmentID:        req.DepartmentID,
	}
	switch req.Gender {
	case api.Female:
		user.Gender = model.GenderFemale
	case api.Male:
		user.Gender = model.GenderMale
	}
	return user
}

func fromAPIPutUserRequest(userID uint64, req api.PutUserJSONRequestBody) model.User {
	user := model.User{
		ShortUserInfo: model.ShortUserInfo{
			ID:                userID,
			Email:             string(req.Email),
			FirstName:         req.FirstName,
			LastName:          req.LastName,
			MiddleName:        req.MiddleName,
			MobilePhoneNumber: req.MobilePhoneNumber,
			OfficePhoneNumber: req.OfficePhoneNumber,
		},
		DateOfBirth:         req.DateOfBirth.Time,
		PlaceOfBirth:        req.PlaceOfBirth,
		Grade:               req.Grade,
		RegistrationAddress: req.RegistrationAddress,
		ResidentialAddress:  req.ResidentialAddress,
		Insurance:           model.Insurance{Number: req.Insurance.Number},
		Taxpayer:            model.Taxpayer{Number: req.Taxpayer.Number},
		PositionID:          req.PositionID,
		DepartmentID:        req.DepartmentID,
	}
	switch req.Gender {
	case api.Female:
		user.Gender = model.GenderFemale
	case api.Male:
		user.Gender = model.GenderMale
	}
	return user
}

func toAPIGetExpandedUserResponse(u *model.ExpandedUser) api.GetExpandedUserResponse {
	var expUser api.GetExpandedUserResponse

	expUser.GetUserResponse = toAPIGetUserResponse(&u.User)
	// 	expUser.Military = nil
	// 	expUser.Educations = ToAPIListEducations(u.Educations)
	// 	expUser.Trainings = ToAPIListTrainings(u.Trainings)
	// 	expUser.Passports = ToAPIPassports(u.Passports)
	// 	expUser.Vacations = ToAPIListVacations(u.Vacations)
	// 	expUser.Contracts = ToAPIListContracts(u.Contracts)

	return expUser
}

func toAPIGetUserResponse(u *model.User) api.GetUserResponse {
	resp := api.GetUserResponse{
		MobilePhoneNumber:   u.MobilePhoneNumber,
		OfficePhoneNumber:   u.OfficePhoneNumber,
		DateOfBirth:         types.Date{Time: u.DateOfBirth},
		DepartmentID:        u.DepartmentID,
		Email:               types.Email(u.Email),
		FirstName:           u.FirstName,
		Gender:              "",
		Grade:               u.Grade,
		ID:                  u.ID,
		LastName:            u.LastName,
		MiddleName:          u.MiddleName,
		PlaceOfBirth:        u.PlaceOfBirth,
		PositionID:          u.PositionID,
		RegistrationAddress: u.RegistrationAddress,
		ResidentialAddress:  u.ResidentialAddress,
		Insurance: api.Insurance{
			Number:  u.Insurance.Number,
			HasScan: &u.Insurance.HasScan,
		},
		Taxpayer: api.Taxpayer{
			Number:  u.Taxpayer.Number,
			HasScan: &u.Taxpayer.HasScan,
		},
		Military: &api.Military{
			Category:    u.Military.Category,
			Comissariat: u.Military.Commissariat,
			HasScan:     &u.Military.HasScan,
			Rank:        u.Military.Rank,
			Speciality:  u.Military.Speciality,
		},
		PersonalDataProcessing: api.PersonalDataProcessing{
			HasScan: u.PersonalDataProcessing.HasScan,
		},
	}
	switch u.Gender {
	case model.GenderFemale:
		resp.Gender = api.Female
	case model.GenderMale:
		resp.Gender = api.Male
	}
	return resp
}

func toAPIListUsers(users []model.ShortUserInfo) []api.ListUsersItem {
	res := make([]api.ListUsersItem, len(users))
	for i := 0; i < len(users); i++ {
		res[i] = toAPIListUser(users[i])
	}
	return res
}

func toAPIListUser(u model.ShortUserInfo) api.ListUsersItem {
	item := api.ListUsersItem{
		ID:                u.ID,
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		MiddleName:        u.MiddleName,
		Email:             types.Email(u.Email),
		Position:          u.Position,
		Department:        u.Department,
		MobilePhoneNumber: u.MobilePhoneNumber,
		OfficePhoneNumber: u.OfficePhoneNumber,
	}
	return item
}
