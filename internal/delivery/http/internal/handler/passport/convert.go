package passport

import (
	"github.com/oapi-codegen/runtime/types"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

func fromAPIAddPassportRequest(req api.AddPassportJSONRequestBody) model.Passport {
	mp := model.Passport{
		Citizenship:  req.Citizenship,
		IssuedBy:     req.IssuedBy,
		IssuedByCode: req.IssuedByCode,
		IssuedDate:   req.IssuedDate.Time,
		Number:       req.Number,
	}
	switch req.Type {
	case api.National:
		mp.Type = model.PassportTypeNational
	case api.International:
		mp.Type = model.PassportTypeInternational
	}
	return mp
}

func fromAPIPutPassportRequest(passportID uint64, req api.PutPassportJSONRequestBody) model.Passport {
	mp := model.Passport{
		Citizenship:  req.Citizenship,
		IssuedBy:     req.IssuedBy,
		IssuedByCode: req.IssuedByCode,
		IssuedDate:   req.IssuedDate.Time,
		Number:       req.Number,
		ID:           passportID,
	}
	switch req.Type {
	case api.International:
		mp.Type = model.PassportTypeInternational
	case api.National:
		mp.Type = model.PassportTypeNational
	}
	return mp
}

func toAPIGetPassportResponse(mp *model.Passport) api.GetPassportResponse {
	resp := api.GetPassportResponse{
		ID:           mp.ID,
		Citizenship:  mp.Citizenship,
		IssuedBy:     mp.IssuedBy,
		IssuedByCode: mp.IssuedByCode,
		IssuedDate:   types.Date{Time: mp.IssuedDate},
		Number:       mp.Number,
	}
	switch mp.Type {
	case model.PassportTypeInternational:
		resp.Type = api.International
	case model.PassportTypeNational:
		resp.Type = api.National
	}
	return resp
}

func toAPIListPassports(eds []model.Passport) api.ListPassportsResponse {
	res := make([]api.Passport, len(eds))
	for i := 0; i < len(eds); i++ {
		res[i] = toAPIPassport(eds[i])
	}
	return res
}

func toAPIPassport(mp model.Passport) api.Passport {
	resp := api.GetPassportResponse{
		ID:           mp.ID,
		Citizenship:  mp.Citizenship,
		IssuedBy:     mp.IssuedBy,
		IssuedByCode: mp.IssuedByCode,
		IssuedDate:   types.Date{Time: mp.IssuedDate},
		Number:       mp.Number,
		HasScan:      mp.HasScan,
	}
	switch mp.Type {
	case model.PassportTypeInternational:
		resp.Type = api.International
	case model.PassportTypeNational:
		resp.Type = api.National
	}
	return resp
}

func toAPIPassports(psps []model.Passport) []api.Passport {
	res := make([]api.Passport, len(psps))
	for i := 0; i < len(psps); i++ {
		res[i] = toAPIPassport(psps[i])
	}
	return res
}
