package handlers

import (
	"net/http"

	"github.com/muonsoft/validation/validator"

	serr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
)

// @Produce application/json
// @Success 200 {array} api.Education
// @Router  /users/{user_id}/educations [get]
func (h *handler) ListEducations(w http.ResponseWriter, r *http.Request, userID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.Education true ""
// @Failure 409  {object} api.Error "education already exists"
// @Router  /users/{user_id}/educations [post]
func (h *handler) AddEducation(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var e api.Education
	// TODO: decode education from request body

	if err := e.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /users/{user_id}/educations/{education_id} [delete]
func (h *handler) DeleteEducation(w http.ResponseWriter, r *http.Request, userID uint64, educationID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Education
// @Router  /users/{user_id}/educations/{education_id} [get]
func (h *handler) GetEducation(w http.ResponseWriter, r *http.Request, userID uint64, educationID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PatchEducationJSONRequestBody true ""
// @Router  /users/{user_id}/educations/{education_id} [patch]
func (h *handler) PatchEducation(w http.ResponseWriter, r *http.Request, userID uint64, educationID uint64) {
	ctx := r.Context()

	var patch api.PatchEducationJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
