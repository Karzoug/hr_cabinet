package convert

import (
	"github.com/oapi-codegen/runtime/types"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

func ToAPIVisas(vs []model.Visa) []api.Visa {
	res := make([]api.Visa, len(vs))
	for i := 0; i < len(vs); i++ {
		res[i] = ToAPIVisa(&vs[i])
	}
	return res
}

func ToAPIVisa(mv *model.Visa) api.Visa {
	var ne api.VisaNumberEntries
	switch mv.NumberEntries {
	case model.VisaNumberEntriesN1:
		ne = api.VisaNumberEntriesN1
	case model.VisaNumberEntriesN2:
		ne = api.VisaNumberEntriesN2
	case model.VisaNumberEntriesMult:
		ne = api.VisaNumberEntriesMult
	}

	return api.Visa{
		ID:            &mv.ID,
		IssuedState:   mv.IssuedState,
		Number:        mv.Number,
		NumberEntries: ne,
		ValidFrom:     types.Date{Time: mv.ValidFrom},
		ValidTo:       types.Date{Time: mv.ValidTo},
	}
}

func ToModelVisa(v api.Visa) model.Visa {
	var ne model.VisaNumberEntries
	switch v.NumberEntries {
	case api.VisaNumberEntriesN1:
		ne = model.VisaNumberEntriesN1
	case api.VisaNumberEntriesN2:
		ne = model.VisaNumberEntriesN2
	case api.VisaNumberEntriesMult:
		ne = model.VisaNumberEntriesMult
	}

	mv := model.Visa{
		Number:        v.Number,
		IssuedState:   v.IssuedState,
		ValidTo:       v.ValidTo.Time,
		ValidFrom:     v.ValidFrom.Time,
		NumberEntries: ne,
	}
	return mv
}
