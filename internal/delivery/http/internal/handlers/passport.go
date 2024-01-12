package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/muonsoft/validation/validator"
	"github.com/oapi-codegen/runtime/types"

	serr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

// @Produce application/json
// @Success 200 {array} api.Passport
// @Router  /users/{user_id}/passports [get]
func (h *handler) ListPassports(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	psps, err := h.userService.ListPassports(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			serr.ErrorMessage(w, r, http.StatusNotFound, user.ErrUserNotFound.Error(), nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}

	if err := response.JSON(w, http.StatusOK, convertPassportsToAPIPassports(psps)); err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
	}
}

// @Accept  application/json
// @Param   body body api.Passport true ""
// @Failure 409  {object} api.Error "passport already exists"
// @Router  /users/{user_id}/passports [post]
func (h *handler) AddPassport(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var p api.Passport
	if err := request.DecodeJSONStrict(w, r, &p); err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := p.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	id, err := h.userService.AddPassport(ctx, userID, convertAPIPassportToPassport(p))
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			serr.ErrorMessage(w, r, http.StatusConflict, user.ErrUserNotFound.Error(), nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+
			"/passports/"+strconv.FormatUint(id, 10))
}

// @Router /users/{user_id}/passports/{passport_id} [delete]
func (h *handler) DeletePassport(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Passport
// @Router  /users/{user_id}/passports/{passport_id} [get]
func (h *handler) GetPassport(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64) {
	ctx := r.Context()

	p, err := h.userService.GetPassport(ctx, userID, passportID)
	if err != nil {
		if errors.Is(err, user.ErrPassportNotFound) {
			serr.ErrorMessage(w, r, http.StatusNotFound, user.ErrPassportNotFound.Error(), nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}

	if err := response.JSON(w, http.StatusOK, convertPassportToAPIPassport(p)); err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
	}
}

// @Accept application/json
// @Param   body body api.PatchPassportJSONRequestBody true ""
// @Router  /users/{user_id}/passports/{passport_id} [patch]
func (h *handler) PatchPassport(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64) {
	ctx := r.Context()

	var patch api.PatchPassportJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

func convertPassportsToAPIPassports(psps []model.Passport) []api.Passport {
	res := make([]api.Passport, len(psps))
	for i := 0; i < len(psps); i++ {
		res[i] = convertPassportToAPIPassport(&psps[i])
	}
	return res
}

func convertPassportToAPIPassport(mp *model.Passport) api.Passport {
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

func convertAPIPassportToPassport(p api.Passport) model.Passport {
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
