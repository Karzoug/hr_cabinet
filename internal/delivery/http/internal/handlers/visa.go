package handlers

import (
	"net/http"
	"strconv"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/convert"
	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
)

// @Produce application/json
// @Success 200 {object} api.ListVisasResponse
// @Router  /users/{user_id}/passports/{passport_id}/visas [get]
func (h *handler) ListVisas(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64) {
	ctx := r.Context()

	vs, err := h.userService.ListVisas(ctx, userID, passportID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIListVisas(vs)); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
	}
}

// @Accept application/json
// @Param   body body api.AddVisaJSONRequestBody true ""
// @Router  /users/{user_id}/passports/{passport_id}/visas [post]
func (h *handler) AddVisa(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64) {
	ctx := r.Context()

	var v api.AddVisaJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &v); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := v.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	id, err := h.userService.AddVisa(ctx, userID, passportID, convert.FromAPIAddVisaRequest(v))
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+
			"/passports/"+strconv.FormatUint(passportID, 10)+
			"/visas/"+strconv.FormatUint(id, 10))
	w.WriteHeader(http.StatusCreated)
}

// @Router /users/{user_id}/passports/{passport_id}/visas/{visa_id} [delete]
func (h *handler) DeleteVisa(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64, visaID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.GetVisaResponse
// @Router  /users/{user_id}/passports/{passport_id}/visas/{visa_id} [get]
func (h *handler) GetVisa(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64, visaID uint64) {
	ctx := r.Context()

	v, err := h.userService.GetVisa(ctx, userID, passportID, visaID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIGetVisaResponse(v)); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
	}
}

// @Accept application/json
// @Param   body body api.PatchVisaJSONRequestBody true ""
// @Router  /users/{user_id}/passports/{passport_id}/visas/{visa_id} [patch]
func (h *handler) PatchVisa(w http.ResponseWriter, r *http.Request, userID, passportID, visaID uint64) {
	ctx := r.Context()

	var patch api.PatchVisaJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PutVisaJSONRequestBody true ""
// @Router  /users/{user_id}/visas/{visa_id} [put]
func (h *handler) PutVisa(w http.ResponseWriter, r *http.Request, userID, passportID, visaID uint64) {
	ctx := r.Context()

	var v api.PutVisaJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &v); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := v.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	err := h.userService.UpdateVisa(ctx, userID, passportID, convert.FromAPIPutVisaRequest(visaID, v))
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}
}
