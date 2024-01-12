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
// @Success 200 {array} api.Visa
// @Router  /users/{user_id}/passports/{passport_id}/visas [get]
func (h *handler) ListVisas(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64) {
	ctx := r.Context()

	vs, err := h.userService.ListVisas(ctx, userID, passportID)
	if err != nil {
		if errors.Is(err, user.ErrUserOrPassportNotFound) {
			serr.ErrorMessage(w, r, http.StatusNotFound, user.ErrUserOrPassportNotFound.Error(), nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}

	if err := response.JSON(w, http.StatusOK, convertVisasToAPIVisas(vs)); err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
	}
}

// @Accept application/json
// @Param   body body api.Visa true ""
// @Failure 409  {object} api.Error "visa already exists"
// @Router  /users/{user_id}/passports/{passport_id}/visas [post]
func (h *handler) AddVisa(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64) {
	ctx := r.Context()

	var v api.Visa
	if err := request.DecodeJSONStrict(w, r, &v); err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := v.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	id, err := h.userService.AddVisa(ctx, userID, passportID, convertAPIVisaToVisa(v))
	if err != nil {
		if errors.Is(err, user.ErrUserOrPassportNotFound) {
			serr.ErrorMessage(w, r, http.StatusConflict, user.ErrUserOrPassportNotFound.Error(), nil)
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
			"/passports/"+strconv.FormatUint(passportID, 10)+
			"/visas/"+strconv.FormatUint(id, 10))
}

// @Router /users/{user_id}/passports/{passport_id}/visas/{visa_id} [delete]
func (h *handler) DeleteVisa(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64, visaID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Visa
// @Router  /users/{user_id}/passports/{passport_id}/visas/{visa_id} [get]
func (h *handler) GetVisa(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64, visaID uint64) {
	ctx := r.Context()

	p, err := h.userService.GetVisa(ctx, userID, passportID, visaID)
	if err != nil {
		if errors.Is(err, user.ErrVisaNotFound) {
			serr.ErrorMessage(w, r, http.StatusNotFound, user.ErrVisaNotFound.Error(), nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}

	if err := response.JSON(w, http.StatusOK, convertVisaToAPIVisa(p)); err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
	}
}

// @Accept application/json
// @Param   body body api.PatchVisaJSONRequestBody true ""
// @Router  /users/{user_id}/passports/{passport_id}/visas/{visa_id} [patch]
func (h *handler) PatchVisa(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64, visaID uint64) {
	ctx := r.Context()

	var patch api.PatchVisaJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

func convertVisasToAPIVisas(vs []model.Visa) []api.Visa {
	res := make([]api.Visa, len(vs))
	for i := 0; i < len(vs); i++ {
		res[i] = convertVisaToAPIVisa(&vs[i])
	}
	return res
}

func convertVisaToAPIVisa(mv *model.Visa) api.Visa {
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

func convertAPIVisaToVisa(v api.Visa) model.Visa {
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
