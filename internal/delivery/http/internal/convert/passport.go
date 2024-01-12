package convert

import (
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/oapi-codegen/runtime/types"
)

func ToAPIPassports(psps []model.Passport) []api.Passport {
	res := make([]api.Passport, len(psps))
	for i := 0; i < len(psps); i++ {
		res[i] = ToAPIPassport(&psps[i])
	}
	return res
}

func ToAPIPassportsWithVisas(psps []model.PassportWithVisas) []api.PassportWithVisas {
	res := make([]api.PassportWithVisas, len(psps))
	for i := 0; i < len(psps); i++ {
		res[i] = ToAPIPassportWithVisas(&psps[i])
	}
	return res
}

func ToAPIPassportWithVisas(mp *model.PassportWithVisas) api.PassportWithVisas {
	return api.PassportWithVisas{
		Passport: ToAPIPassport(&mp.Passport),
		Visas:    ToAPIVisas(mp.Visas),
	}
}

func ToAPIPassport(mp *model.Passport) api.Passport {
	var pt api.PassportType
	switch mp.Type {
	case model.PassportTypeInternal:
		pt = api.PassportTypeInternal
	case model.PassportTypeExternal:
		pt = api.PassportTypeExternal
	case model.PassportTypeForeigners:
		pt = api.PassportTypeForeigners
	}

	return api.Passport{
		ID:         &mp.ID,
		IssuedBy:   mp.IssuedBy,
		IssuedDate: types.Date{Time: mp.IssuedDate},
		Number:     mp.Number,
		Type:       pt,
		VisasCount: mp.VisasCount,
	}
}

func ToModelPassport(p api.Passport) model.Passport {
	var pt model.PassportType
	switch p.Type {
	case api.PassportTypeInternal:
		pt = model.PassportTypeInternal
	case api.PassportTypeExternal:
		pt = model.PassportTypeExternal
	case api.PassportTypeForeigners:
		pt = model.PassportTypeForeigners
	}

	mp := model.Passport{
		IssuedBy:   p.IssuedBy,
		IssuedDate: p.IssuedDate.Time,
		Number:     p.Number,
		Type:       pt,
	}
	return mp
}
