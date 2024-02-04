package vacation

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

type vacationHandlers struct {
	usecase user.VacationUseCase
	logger  *slog.Logger
}

func newVacationHandlers(v user.VacationUseCase, l *slog.Logger) vacationHandlers {
	return vacationHandlers{
		usecase: v,
		logger:  l,
	}
}

// @Produce application/json
// @Success 200 {object} api.ListVacationsResponse
// @Router  /users/{user_id}/vacations [get]
func (h vacationHandlers) ListVacations(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	vacations, err := h.usecase.List(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, toAPIListVacations(vacations)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

// @Accept application/json
// @Param   body body api.AddVacationJSONRequestBody true ""
// @Failure 409  {object} api.Error "vacation already exists"
// @Router  /users/{user_id}/vacations [post]
func (h vacationHandlers) AddVacation(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var v api.AddVacationJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &v); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := v.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	id, err := h.usecase.Add(ctx, userID, fromAPIAddVacationRequest(v))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+
			"/vacations/"+strconv.FormatUint(id, 10))
	w.WriteHeader(http.StatusCreated)
}

// @Router /users/{user_id}/vacations/{vacation_id} [delete]
func (h vacationHandlers) DeleteVacation(w http.ResponseWriter, r *http.Request, userID uint64, vacationID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept application/json
// @Success 200 {object} api.GetVacationResponse
// @Router  /users/{user_id}/vacations/{vacation_id} [get]
func (h vacationHandlers) GetVacation(w http.ResponseWriter, r *http.Request, userID uint64, vacationID uint64) {
	ctx := r.Context()

	v, err := h.usecase.Get(ctx, userID, vacationID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, toAPIGetVacationResponse(v)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

// @Accept application/json
// @Param   body body api.PatchVacationJSONRequestBody true ""
// @Router  /users/{user_id}/vacations/{vacation_id} [patch]
func (h vacationHandlers) PatchVacation(w http.ResponseWriter, r *http.Request, userID uint64, vacationID uint64) {
	ctx := r.Context()

	var patch api.PatchVacationJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PutVacationJSONRequestBody true ""
// @Router  /users/{user_id}/vacations/{vacation_id} [put]
func (h vacationHandlers) PutVacation(w http.ResponseWriter, r *http.Request, userID, vacationID uint64) {
	ctx := r.Context()

	var v api.PutVacationJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &v); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := v.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	err := h.usecase.Update(ctx, userID, fromAPIPutVacationRequest(vacationID, v))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}
}
