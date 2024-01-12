package convert

import (
	"github.com/oapi-codegen/runtime/types"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

func ToAPITrainings(eds []model.Training) []api.Training {
	res := make([]api.Training, len(eds))
	for i := 0; i < len(eds); i++ {
		res[i] = ToAPITraining(&eds[i])
	}
	return res
}

func ToAPITraining(mtr *model.Training) api.Training {
	return api.Training{
		Cost:              (api.Money)(mtr.Cost),
		DateFrom:          types.Date{Time: mtr.DateFrom},
		DateTo:            types.Date{Time: mtr.DateTo},
		ID:                &mtr.ID,
		IssuedInstitution: mtr.IssuedInstitution,
		Program:           mtr.Program,
	}
}

func ToModelTraining(tr api.Training) model.Training {
	return model.Training{
		Cost:              (uint64)(tr.Cost),
		DateFrom:          tr.DateFrom.Time,
		DateTo:            tr.DateTo.Time,
		IssuedInstitution: tr.IssuedInstitution,
		Program:           tr.Program,
	}
}
