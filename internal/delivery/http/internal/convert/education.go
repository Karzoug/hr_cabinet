package convert

import (
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/oapi-codegen/runtime/types"
)

func ToAPIEducations(eds []model.Education) []api.Education {
	res := make([]api.Education, len(eds))
	for i := 0; i < len(eds); i++ {
		res[i] = ToAPIEducation(&eds[i])
	}
	return res
}

func ToAPIEducation(med *model.Education) api.Education {
	return api.Education{
		DateFrom:          types.Date{Time: med.DateFrom},
		DateTo:            types.Date{Time: med.DateTo},
		ID:                &med.ID,
		IssuedInstitution: med.IssuedInstitution,
		Number:            med.Number,
		Program:           med.Program,
	}
}

func ToModelEducation(e api.Education) model.Education {
	med := model.Education{
		DateFrom:          e.DateFrom.Time,
		DateTo:            e.DateTo.Time,
		IssuedInstitution: e.IssuedInstitution,
		Number:            e.Number,
		Program:           e.Program,
	}
	return med
}
