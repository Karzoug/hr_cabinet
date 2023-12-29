package handlers

import (
	"net/http"

	"github.com/muonsoft/validation/validator"

	serr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
)

// @Produce application/json
// @Success 200 {array} api.Visa
// @Router  /users/{user_id}/passports/{passport_id}/visas [get]
func (h *handler) ListVisas(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept application/json
// @Param   body body api.Visa true ""
// @Failure 409  {object} api.Error "visa already exists"
// @Router  /users/{user_id}/passports/{passport_id}/visas [post]
func (h *handler) AddVisa(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64) {
	ctx := r.Context()

	var v api.Visa
	// TODO: decode visa from request body
	if err := v.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /users/{user_id}/passports/{passport_id}/visas/{visa_id} [delete]
func (h *handler) DeleteVisa(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64, visaID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Visa
// @Router  /users/{user_id}/passports/{passport_id}/visas/{visa_id} [get]
func (h *handler) GetVisa(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64, visaID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
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
