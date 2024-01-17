package convert

import (
	"github.com/oapi-codegen/runtime/types"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

func FromAPIAddVisaRequest(req api.AddVisaJSONRequestBody) model.Visa {
	v := model.Visa{
		Number:        req.Number,
		IssuedState:   req.IssuedState,
		ValidTo:       req.ValidTo.Time,
		ValidFrom:     req.ValidFrom.Time,
		NumberEntries: model.VisaNumberEntries(req.NumberEntries),
	}
	switch req.NumberEntries {
	case api.N1:
		v.NumberEntries = model.VisaNumberEntriesN1
	case api.N2:
		v.NumberEntries = model.VisaNumberEntriesN2
	case api.Mult:
		v.NumberEntries = model.VisaNumberEntriesMult
	}
	return v
}

func FromAPIPutVisaRequest(visaID uint64, req api.PutVisaJSONRequestBody) model.Visa {
	v := model.Visa{
		ID:            visaID,
		Number:        req.Number,
		IssuedState:   req.IssuedState,
		ValidTo:       req.ValidTo.Time,
		ValidFrom:     req.ValidFrom.Time,
		NumberEntries: model.VisaNumberEntries(req.NumberEntries),
	}
	switch req.NumberEntries {
	case api.N1:
		v.NumberEntries = model.VisaNumberEntriesN1
	case api.N2:
		v.NumberEntries = model.VisaNumberEntriesN2
	case api.Mult:
		v.NumberEntries = model.VisaNumberEntriesMult
	}
	return v
}

func ToAPIGetVisaResponse(mv *model.Visa) api.GetVisaResponse {
	return api.GetVisaResponse{
		ID:            mv.ID,
		Number:        mv.Number,
		IssuedState:   mv.IssuedState,
		ValidTo:       types.Date{Time: mv.ValidTo},
		ValidFrom:     types.Date{Time: mv.ValidFrom},
		NumberEntries: api.VisaNumberEntries(mv.NumberEntries),
	}
}

func ToAPIListVisas(vs []model.Visa) api.ListVisasResponse {
	return toAPIVisas(vs)
}

func toAPIVisas(vs []model.Visa) []api.Visa {
	res := make([]api.Visa, len(vs))
	for i := 0; i < len(vs); i++ {
		res[i] = toAPIVisa(vs[i])
	}
	return res
}

func toAPIVisa(mv model.Visa) api.Visa {
	v := api.Visa{
		ID:          mv.ID,
		IssuedState: mv.IssuedState,
		Number:      mv.Number,
		ValidFrom:   types.Date{Time: mv.ValidFrom},
		ValidTo:     types.Date{Time: mv.ValidTo},
	}
	switch mv.NumberEntries {
	case model.VisaNumberEntriesN1:
		v.NumberEntries = api.N1
	case model.VisaNumberEntriesN2:
		v.NumberEntries = api.N2
	case model.VisaNumberEntriesMult:
		v.NumberEntries = api.Mult
	}
	return v
}
