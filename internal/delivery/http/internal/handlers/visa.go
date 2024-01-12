package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/muonsoft/validation/validator"

	serr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/convert"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user"
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

	if err := response.JSON(w, http.StatusOK, convert.ToAPIVisas(vs)); err != nil {
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

	id, err := h.userService.AddVisa(ctx, userID, passportID, convert.ToModelVisa(v))
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

	if err := response.JSON(w, http.StatusOK, convert.ToAPIVisa(p)); err != nil {
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
