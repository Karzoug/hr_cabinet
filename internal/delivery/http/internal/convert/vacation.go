package convert

import (
	"github.com/oapi-codegen/runtime/types"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

func ToAPIVacations(vcs []model.Vacation) []api.Vacation {
	res := make([]api.Vacation, len(vcs))
	for i := 0; i < len(vcs); i++ {
		res[i] = ToAPIVacation(&vcs[i])
	}
	return res
}

func ToAPIVacation(mv *model.Vacation) api.Vacation {
	return api.Vacation{
		ID:       &mv.ID,
		DateFrom: types.Date{Time: mv.DateBegin},
		DateTo:   types.Date{Time: mv.DateEnd},
	}
}

func ToModelVacation(tr api.Vacation) model.Vacation {
	return model.Vacation{
		DateBegin: tr.DateFrom.Time,
		DateEnd:   tr.DateTo.Time,
	}
}
