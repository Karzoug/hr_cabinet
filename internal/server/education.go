package server

import (
	"net/http"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
)

// @Produce application/json
// @Success 200 {array} api.Education
// @Router  /users/{user_id}/educations [get]
func (s *server) ListEducations(w http.ResponseWriter, r *http.Request, userID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.Education true ""
// @Failure 409  {object} api.Error "education already exists"
// @Router  /users/{user_id}/educations [post]
func (s *server) AddEducation(w http.ResponseWriter, r *http.Request, userID int) {
	ctx := r.Context()

	var e api.Education
	// TODO: decode education from request body

	if err := e.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = api.NewBadRequestErrorFromError(err)
		// encode error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /users/{user_id}/educations/{education_id} [delete]
func (s *server) DeleteEducation(w http.ResponseWriter, r *http.Request, userID int, educationID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Education
// @Router  /users/{user_id}/educations/{education_id} [get]
func (s *server) GetEducation(w http.ResponseWriter, r *http.Request, userID int, educationID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PatchEducationJSONRequestBody true ""
// @Router  /users/{user_id}/educations/{education_id} [patch]
func (s *server) PatchEducation(w http.ResponseWriter, r *http.Request, userID int, educationID int) {
	ctx := r.Context()

	var patch api.PatchEducationJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = api.NewBadRequestErrorFromError(err)
		// encode error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
