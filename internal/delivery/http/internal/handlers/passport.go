package handlers

import (
	"net/http"

	"github.com/muonsoft/validation/validator"

	serr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
)

// @Produce application/json
// @Success 200 {array} api.Passport
// @Router  /users/{user_id}/passports [get]
func (h *handler) ListPassports(w http.ResponseWriter, r *http.Request, userID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.Passport true ""
// @Failure 409  {object} api.Error "passport already exists"
// @Router  /users/{user_id}/passports [post]
func (h *handler) AddPassport(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var p api.Passport
	// TODO: decode passport from request body

	if err := p.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /users/{user_id}/passports/{passport_id} [delete]
func (h *handler) DeletePassport(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Passport
// @Router  /users/{user_id}/passports/{passport_id} [get]
func (h *handler) GetPassport(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
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
