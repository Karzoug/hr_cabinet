package convert

import (
	"github.com/oapi-codegen/runtime/types"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

func FromAPIAddTrainingRequest(req api.AddTrainingJSONRequestBody) model.Training {
	return model.Training{
		Cost:              req.Cost,
		DateFrom:          req.DateFrom.Time,
		DateTo:            req.DateTo.Time,
		IssuedInstitution: req.IssuedInstitution,
		Program:           req.Program,
	}
}

func FromAPIPutTrainingRequest(trainingID uint64, req api.PutTrainingJSONRequestBody) model.Training {
	return model.Training{
		ID:                trainingID,
		Cost:              req.Cost,
		DateFrom:          req.DateFrom.Time,
		DateTo:            req.DateTo.Time,
		IssuedInstitution: req.IssuedInstitution,
		Program:           req.Program,
	}
}

func ToAPIGetTrainingResponse(mtr *model.Training) api.GetTrainingResponse {
	return api.GetTrainingResponse{
		Cost:              mtr.Cost,
		DateFrom:          types.Date{Time: mtr.DateFrom},
		DateTo:            types.Date{Time: mtr.DateTo},
		ID:                mtr.ID,
		IssuedInstitution: mtr.IssuedInstitution,
		Program:           mtr.Program,
	}
}

func ToAPIListTrainings(trs []model.Training) api.ListTrainingsResponse {
	res := make([]api.Training, len(trs))
	for i := 0; i < len(trs); i++ {
		res[i] = toAPITraining(trs[i])
	}
	return res
}

func toAPITraining(mtr model.Training) api.Training {
	return api.Training{
		Cost:              mtr.Cost,
		DateFrom:          types.Date{Time: mtr.DateFrom},
		DateTo:            types.Date{Time: mtr.DateTo},
		ID:                mtr.ID,
		IssuedInstitution: mtr.IssuedInstitution,
		Program:           mtr.Program,
	}
}
