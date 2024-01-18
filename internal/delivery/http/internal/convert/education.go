package convert

import (
	"github.com/oapi-codegen/runtime/types"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

func FromAPIAddEducationRequest(req api.AddEducationJSONRequestBody) model.Education {
	return model.Education{
		DateFrom:          req.DateFrom.Time,
		DateTo:            req.DateTo.Time,
		IssuedInstitution: req.IssuedInstitution,
		Number:            req.Number,
		Program:           req.Program,
	}
}

func FromAPIPutEducationRequest(educationID uint64, req api.PutEducationJSONRequestBody) model.Education {
	return model.Education{
		ID:                educationID,
		DateFrom:          req.DateFrom.Time,
		DateTo:            req.DateTo.Time,
		IssuedInstitution: req.IssuedInstitution,
		Number:            req.Number,
		Program:           req.Program,
	}
}

func ToAPIGetEducationResponse(med *model.Education) api.GetEducationResponse {
	return api.GetEducationResponse{
		ID:                med.ID,
		DateFrom:          types.Date{Time: med.DateFrom},
		DateTo:            types.Date{Time: med.DateTo},
		IssuedInstitution: med.IssuedInstitution,
		Number:            med.Number,
		Program:           med.Program,
		HasScan:           &med.HasScan,
	}
}

func ToAPIListEducations(eds []model.Education) api.ListEducationsResponse {
	res := make([]api.Education, len(eds))
	for i := 0; i < len(eds); i++ {
		res[i] = toAPIEducation(eds[i])
	}
	return res
}

func toAPIEducation(med model.Education) api.Education {
	return api.Education{
		DateFrom:          types.Date{Time: med.DateFrom},
		DateTo:            types.Date{Time: med.DateTo},
		ID:                med.ID,
		IssuedInstitution: med.IssuedInstitution,
		Number:            med.Number,
		Program:           med.Program,
		HasScan:           &med.HasScan,
	}
}
