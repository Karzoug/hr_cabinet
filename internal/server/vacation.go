package server

import (
	"net/http"

	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
	"github.com/muonsoft/validation/validator"
)

// @Produce application/json
// @Success 200 {array} api.Vacation
// @Router  /users/{user_id}/vacations [get]
func (s *server) ListVacations(w http.ResponseWriter, r *http.Request, userID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept application/json
// @Param   body body api.Vacation true ""
// @Failure 409  {object} api.Error "vacation already exists"
// @Router  /users/{user_id}/vacations [post]
func (s *server) AddVacation(w http.ResponseWriter, r *http.Request, userID int) {
	ctx := r.Context()

	var v api.Vacation
	// TODO: decode vacation from request body
	if err := v.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = api.NewBadRequestErrorFromError(err)
		// encode error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /users/{user_id}/vacations/{vacation_id} [delete]
func (s *server) DeleteVacation(w http.ResponseWriter, r *http.Request, userID int, vacationID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept application/json
// @Success 200 {object} api.Vacation
// @Router  /users/{user_id}/vacations/{vacation_id} [get]
func (s *server) GetVacation(w http.ResponseWriter, r *http.Request, userID int, vacationID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept application/json
// @Param   body body api.PatchVacationJSONRequestBody true ""
// @Router  /users/{user_id}/vacations/{vacation_id} [patch]
func (s *server) PatchVacation(w http.ResponseWriter, r *http.Request, userID int, vacationID int) {
	ctx := r.Context()

	var patch api.PatchVacationJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = api.NewBadRequestErrorFromError(err)
		// encode error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
