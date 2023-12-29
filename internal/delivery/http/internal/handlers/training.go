package handlers

import (
	"net/http"

	"github.com/muonsoft/validation/validator"

	serr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
)

// @Produce application/json
// @Success 200 {array} api.Training
// @Router  /users/{user_id}/trainings [get]
func (h *handler) ListTrainings(w http.ResponseWriter, r *http.Request, userID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.Training true ""
// @Failure 409  {object} api.Error "training already exists"
// @Router  /users/{user_id}/trainings [post]
func (h *handler) AddTraining(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var t api.Training
	// TODO: decode training from request body

	if err := t.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /users/{user_id}/trainings/{training_id} [delete]
func (h *handler) DeleteTraining(w http.ResponseWriter, r *http.Request, userID uint64, trainingID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Training
// @Router  /users/{user_id}/trainings/{training_id} [get]
func (h *handler) GetTraining(w http.ResponseWriter, r *http.Request, userID uint64, trainingID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept application/json
// @Param   body body api.PatchTrainingJSONRequestBody true ""
// @Router  /users/{user_id}/trainings/{training_id} [patch]
func (h *handler) PatchTraining(w http.ResponseWriter, r *http.Request, userID uint64, trainingID uint64) {
	ctx := r.Context()

	var patch api.PatchTrainingJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
