package vacation

import (
	"github.com/oapi-codegen/runtime/types"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

func fromAPIAddVacationRequest(req api.AddVacationJSONRequestBody) model.Vacation {
	return model.Vacation{
		DateBegin: req.DateFrom.Time,
		DateEnd:   req.DateTo.Time,
	}
}

func fromAPIPutVacationRequest(vacationID uint64, req api.PutVacationJSONRequestBody) model.Vacation {
	return model.Vacation{
		ID:        vacationID,
		DateBegin: req.DateFrom.Time,
		DateEnd:   req.DateTo.Time,
	}
}

func toAPIGetVacationResponse(med *model.Vacation) api.GetVacationResponse {
	return api.GetVacationResponse{
		ID:       med.ID,
		DateFrom: types.Date{Time: med.DateBegin},
		DateTo:   types.Date{Time: med.DateEnd},
	}
}

func toAPIListVacations(eds []model.Vacation) api.ListVacationsResponse {
	res := make([]api.Vacation, len(eds))
	for i := 0; i < len(eds); i++ {
		res[i] = toAPIVacation(eds[i])
	}
	return res
}

func toAPIVacation(med model.Vacation) api.Vacation {
	return api.Vacation{
		ID:       med.ID,
		DateFrom: types.Date{Time: med.DateBegin},
		DateTo:   types.Date{Time: med.DateEnd},
	}
}
