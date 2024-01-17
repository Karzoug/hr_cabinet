package handlers

import (
	"net/http"
	"strconv"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/convert"
	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
)

// @Produce application/json
// @Success 200 {object} api.ListTrainingsResponse
// @Router  /users/{user_id}/trainings [get]
func (h *handler) ListTrainings(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	trs, err := h.userService.ListTrainings(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIListTrainings(trs)); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
	}
}

// @Accept  application/json
// @Param   body body api.AddTrainingJSONRequestBody true ""
// @Failure 409  {object} api.Error "training already exists"
// @Router  /users/{user_id}/trainings [post]
func (h *handler) AddTraining(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var tr api.AddTrainingJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &tr); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := tr.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	id, err := h.userService.AddTraining(ctx, userID, convert.FromAPIAddTrainingRequest(tr))
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+
			"/trainings/"+strconv.FormatUint(id, 10))
	w.WriteHeader(http.StatusCreated)
}

// @Router /users/{user_id}/trainings/{training_id} [delete]
func (h *handler) DeleteTraining(w http.ResponseWriter, r *http.Request, userID uint64, trainingID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.GetTrainingResponse
// @Router  /users/{user_id}/trainings/{training_id} [get]
func (h *handler) GetTraining(w http.ResponseWriter, r *http.Request, userID uint64, trainingID uint64) {
	ctx := r.Context()

	tr, err := h.userService.GetTraining(ctx, userID, trainingID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIGetTrainingResponse(tr)); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
	}
}

// @Accept application/json
// @Param   body body api.PatchTrainingJSONRequestBody true ""
// @Router  /users/{user_id}/trainings/{training_id} [patch]
func (h *handler) PatchTraining(w http.ResponseWriter, r *http.Request, userID uint64, trainingID uint64) {
	ctx := r.Context()

	var patch api.PatchTrainingJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &patch); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PutTrainingJSONRequestBody true ""
// @Router  /users/{user_id}/trainings/{training_id} [put]
func (h *handler) PutTraining(w http.ResponseWriter, r *http.Request, userID, trainingID uint64) {
	ctx := r.Context()

	var tr api.PutTrainingJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &tr); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := tr.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	err := h.userService.UpdateTraining(ctx, userID, convert.FromAPIPutTrainingRequest(trainingID, tr))
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}
}
