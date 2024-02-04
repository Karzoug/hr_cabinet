package visa

import (
	"github.com/oapi-codegen/runtime/types"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

func fromAPIAddVisaRequest(req api.AddVisaJSONRequestBody) model.Visa {
	v := model.Visa{
		Number:      req.Number,
		IssuedState: req.IssuedState,
		ValidTo:     req.ValidTo.Time,
		ValidFrom:   req.ValidFrom.Time,
		Type:        model.VisaType(req.Type),
	}
	return v
}

func fromAPIPutVisaRequest(visaID uint64, req api.PutVisaJSONRequestBody) model.Visa {
	v := model.Visa{
		ID:          visaID,
		Number:      req.Number,
		IssuedState: req.IssuedState,
		ValidTo:     req.ValidTo.Time,
		ValidFrom:   req.ValidFrom.Time,
		Type:        model.VisaType(req.Type),
	}
	return v
}

func toAPIGetVisaResponse(mv *model.Visa) api.GetVisaResponse {
	return api.GetVisaResponse{
		ID:          mv.ID,
		Number:      mv.Number,
		IssuedState: mv.IssuedState,
		ValidTo:     types.Date{Time: mv.ValidTo},
		ValidFrom:   types.Date{Time: mv.ValidFrom},
		Type:        string(mv.Type),
	}
}

func toAPIListVisas(vs []model.Visa) api.ListVisasResponse {
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
		Type:        string(mv.Type),
	}
	return v
}
