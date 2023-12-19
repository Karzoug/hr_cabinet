package handlers

import (
	"net/http"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
)

// @Produce application/json
// @Success 200 {array} api.Passport
// @Router  /users/{user_id}/passports [get]
func (s *server) ListPassports(w http.ResponseWriter, r *http.Request, userID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.Passport true ""
// @Failure 409  {object} api.Error "passport already exists"
// @Router  /users/{user_id}/passports [post]
func (s *server) AddPassport(w http.ResponseWriter, r *http.Request, userID int) {
	ctx := r.Context()

	var p api.Passport
	// TODO: decode passport from request body

	if err := p.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = api.NewBadRequestErrorFromError(err)
		// encode error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /users/{user_id}/passports/{passport_id} [delete]
func (s *server) DeletePassport(w http.ResponseWriter, r *http.Request, userID int, passportID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Passport
// @Router  /users/{user_id}/passports/{passport_id} [get]
func (s *server) GetPassport(w http.ResponseWriter, r *http.Request, userID int, passportID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept application/json
// @Param   body body api.PatchPassportJSONRequestBody true ""
// @Router  /users/{user_id}/passports/{passport_id} [patch]
func (s *server) PatchPassport(w http.ResponseWriter, r *http.Request, userID int, passportID int) {
	ctx := r.Context()

	var patch api.PatchPassportJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = api.NewBadRequestErrorFromError(err)
		// encode error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
