package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/muonsoft/validation/validator"
	"github.com/oapi-codegen/runtime/types"

	serr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

// @Produce application/json
// @Success 200 {array} api.Training
// @Router  /users/{user_id}/trainings [get]
func (h *handler) ListTrainings(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	eds, err := h.userService.ListTrainings(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			serr.ErrorMessage(w, r, http.StatusNotFound, user.ErrUserNotFound.Error(), nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}

	if err := response.JSON(w, http.StatusOK, convertTrainingsToAPITrainings(eds)); err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
	}
}

// @Accept  application/json
// @Param   body body api.Training true ""
// @Failure 409  {object} api.Error "training already exists"
// @Router  /users/{user_id}/trainings [post]
func (h *handler) AddTraining(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var tr api.Training
	if err := request.DecodeJSONStrict(w, r, &tr); err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := tr.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	id, err := h.userService.AddTraining(ctx, userID, convertAPITrainingToTraining(tr))
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			serr.ErrorMessage(w, r, http.StatusConflict, user.ErrUserNotFound.Error(), nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+
			"/trainings/"+strconv.FormatUint(id, 10))
}

// @Router /users/{user_id}/trainings/{training_id} [delete]
func (h *handler) DeleteTraining(w http.ResponseWriter, r *http.Request, userID uint64, trainingID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Training
// @Router  /users/{user_id}/trainings/{training_id} [get]
func (h *handler) GetTraining(w http.ResponseWriter, r *http.Request, userID uint64, trainingID uint64) {
	ctx := r.Context()

	ed, err := h.userService.GetTraining(ctx, userID, trainingID)
	if err != nil {
		if errors.Is(err, user.ErrTrainingNotFound) {
			serr.ErrorMessage(w, r, http.StatusNotFound, user.ErrTrainingNotFound.Error(), nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}

	if err := response.JSON(w, http.StatusOK, convertTrainingToAPITraining(ed)); err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
	}
}

// @Accept application/json
// @Param   body body api.PatchTrainingJSONRequestBody true ""
// @Router  /users/{user_id}/trainings/{training_id} [patch]
func (h *handler) PatchTraining(w http.ResponseWriter, r *http.Request, userID uint64, trainingID uint64) {
	ctx := r.Context()

	var patch api.PatchTrainingJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &patch); err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

func convertTrainingsToAPITrainings(eds []model.Training) []api.Training {
	res := make([]api.Training, len(eds))
	for i := 0; i < len(eds); i++ {
		res[i] = convertTrainingToAPITraining(&eds[i])
	}
	return res
}

func convertTrainingToAPITraining(mtr *model.Training) api.Training {
	return api.Training{
		Cost:              (api.Money)(mtr.Cost),
		DateFrom:          types.Date{Time: mtr.DateFrom},
		DateTo:            types.Date{Time: mtr.DateTo},
		ID:                &mtr.ID,
		IssuedInstitution: mtr.IssuedInstitution,
		Program:           mtr.Program,
	}
}

func convertAPITrainingToTraining(tr api.Training) model.Training {
	return model.Training{
		Cost:              (uint64)(tr.Cost),
		DateFrom:          tr.DateFrom.Time,
		DateTo:            tr.DateTo.Time,
		IssuedInstitution: tr.IssuedInstitution,
		Program:           tr.Program,
	}
}
