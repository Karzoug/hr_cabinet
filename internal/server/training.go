package server

import (
	"net/http"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
)

// @Produce application/json
// @Success 200 {array} api.Training
// @Router  /users/{user_id}/trainings [get]
func (s *server) ListTrainings(w http.ResponseWriter, r *http.Request, userID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.Training true ""
// @Failure 409  {object} api.Error "training already exists"
// @Router  /users/{user_id}/trainings [post]
func (s *server) AddTraining(w http.ResponseWriter, r *http.Request, userID int) {
	ctx := r.Context()

	var t api.Training
	// TODO: decode training from request body

	if err := t.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = api.NewBadRequestErrorFromError(err)
		// encode error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /users/{user_id}/trainings/{training_id} [delete]
func (s *server) DeleteTraining(w http.ResponseWriter, r *http.Request, userID int, trainingID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Training
// @Router  /users/{user_id}/trainings/{training_id} [get]
func (s *server) GetTraining(w http.ResponseWriter, r *http.Request, userID int, trainingID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept application/json
// @Param   body body api.PatchTrainingJSONRequestBody true ""
// @Router  /users/{user_id}/trainings/{training_id} [patch]
func (s *server) PatchTraining(w http.ResponseWriter, r *http.Request, userID int, trainingID int) {
	ctx := r.Context()

	var patch api.PatchTrainingJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = api.NewBadRequestErrorFromError(err)
		// encode error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
