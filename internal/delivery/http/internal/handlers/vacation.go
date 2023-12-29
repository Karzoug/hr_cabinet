package handlers

import (
	"net/http"

	"github.com/muonsoft/validation/validator"

	serr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
)

// @Produce application/json
// @Success 200 {array} api.Vacation
// @Router  /users/{user_id}/vacations [get]
func (h *handler) ListVacations(w http.ResponseWriter, r *http.Request, userID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept application/json
// @Param   body body api.Vacation true ""
// @Failure 409  {object} api.Error "vacation already exists"
// @Router  /users/{user_id}/vacations [post]
func (h *handler) AddVacation(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var v api.Vacation
	// TODO: decode vacation from request body
	if err := v.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /users/{user_id}/vacations/{vacation_id} [delete]
func (h *handler) DeleteVacation(w http.ResponseWriter, r *http.Request, userID uint64, vacationID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept application/json
// @Success 200 {object} api.Vacation
// @Router  /users/{user_id}/vacations/{vacation_id} [get]
func (h *handler) GetVacation(w http.ResponseWriter, r *http.Request, userID uint64, vacationID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept application/json
// @Param   body body api.PatchVacationJSONRequestBody true ""
// @Router  /users/{user_id}/vacations/{vacation_id} [patch]
func (h *handler) PatchVacation(w http.ResponseWriter, r *http.Request, userID uint64, vacationID uint64) {
	ctx := r.Context()

	var patch api.PatchVacationJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
