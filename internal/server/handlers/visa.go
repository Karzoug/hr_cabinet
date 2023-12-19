package handlers

import (
	"net/http"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
)

// @Produce application/json
// @Success 200 {array} api.Visa
// @Router  /users/{user_id}/passports/{passport_id}/visas [get]
func (s *server) ListVisas(w http.ResponseWriter, r *http.Request, userID int, passportID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept application/json
// @Param   body body api.Visa true ""
// @Failure 409  {object} api.Error "visa already exists"
// @Router  /users/{user_id}/passports/{passport_id}/visas [post]
func (s *server) AddVisa(w http.ResponseWriter, r *http.Request, userID int, passportID int) {
	ctx := r.Context()

	var v api.Visa
	// TODO: decode visa from request body
	if err := v.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = api.NewBadRequestErrorFromError(err)
		// encode error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /users/{user_id}/passports/{passport_id}/visas/{visa_id} [delete]
func (s *server) DeleteVisa(w http.ResponseWriter, r *http.Request, userID int, passportID int, visaID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Visa
// @Router  /users/{user_id}/passports/{passport_id}/visas/{visa_id} [get]
func (s *server) GetVisa(w http.ResponseWriter, r *http.Request, userID int, passportID int, visaID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept application/json
// @Param   body body api.PatchVisaJSONRequestBody true ""
// @Router  /users/{user_id}/passports/{passport_id}/visas/{visa_id} [patch]
func (s *server) PatchVisa(w http.ResponseWriter, r *http.Request, userID int, passportID int, visaID int) {
	ctx := r.Context()

	var patch api.PatchVisaJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = api.NewBadRequestErrorFromError(err)
		// encode error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
