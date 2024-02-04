package training

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user"
)

type trainingHandlers struct {
	usecase user.TrainingUseCase
	logger  *slog.Logger
}

func newTrainingHandlers(t user.TrainingUseCase, l *slog.Logger) trainingHandlers {
	return trainingHandlers{
		usecase: t,
		logger:  l,
	}
}

// @Produce application/json
// @Success 200 {object} api.ListTrainingsResponse
// @Router  /users/{user_id}/trainings [get]
func (h trainingHandlers) ListTrainings(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	trs, err := h.usecase.List(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, toAPIListTrainings(trs)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

// @Accept  application/json
// @Param   body body api.AddTrainingJSONRequestBody true ""
// @Failure 409  {object} api.Error "training already exists"
// @Router  /users/{user_id}/trainings [post]
func (h trainingHandlers) AddTraining(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var tr api.AddTrainingJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &tr); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := tr.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	id, err := h.usecase.Add(ctx, userID, fromAPIAddTrainingRequest(tr))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+
			"/trainings/"+strconv.FormatUint(id, 10))
	w.WriteHeader(http.StatusCreated)
}

// @Router /users/{user_id}/trainings/{training_id} [delete]
func (h trainingHandlers) DeleteTraining(w http.ResponseWriter, r *http.Request, userID uint64, trainingID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.GetTrainingResponse
// @Router  /users/{user_id}/trainings/{training_id} [get]
func (h trainingHandlers) GetTraining(w http.ResponseWriter, r *http.Request, userID uint64, trainingID uint64) {
	ctx := r.Context()

	tr, err := h.usecase.Get(ctx, userID, trainingID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, toAPIGetTrainingResponse(tr)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

// @Accept application/json
// @Param   body body api.PatchTrainingJSONRequestBody true ""
// @Router  /users/{user_id}/trainings/{training_id} [patch]
func (h trainingHandlers) PatchTraining(w http.ResponseWriter, r *http.Request, userID uint64, trainingID uint64) {
	ctx := r.Context()

	var patch api.PatchTrainingJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &patch); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PutTrainingJSONRequestBody true ""
// @Router  /users/{user_id}/trainings/{training_id} [put]
func (h trainingHandlers) PutTraining(w http.ResponseWriter, r *http.Request, userID, trainingID uint64) {
	ctx := r.Context()

	var tr api.PutTrainingJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &tr); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := tr.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	err := h.usecase.Update(ctx, userID, fromAPIPutTrainingRequest(trainingID, tr))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}
}
