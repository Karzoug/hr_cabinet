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
// @Success 200 {object} api.ListVacationsResponse
// @Router  /users/{user_id}/vacations [get]
func (h *handler) ListVacations(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	vacations, err := h.userService.ListVacations(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIListVacations(vacations)); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
	}
}

// @Accept application/json
// @Param   body body api.AddVacationJSONRequestBody true ""
// @Failure 409  {object} api.Error "vacation already exists"
// @Router  /users/{user_id}/vacations [post]
func (h *handler) AddVacation(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var v api.AddVacationJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &v); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := v.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	id, err := h.userService.AddVacation(ctx, userID, convert.FromAPIAddVacationRequest(v))
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+
			"/vacations/"+strconv.FormatUint(id, 10))
	w.WriteHeader(http.StatusCreated)
}

// @Router /users/{user_id}/vacations/{vacation_id} [delete]
func (h *handler) DeleteVacation(w http.ResponseWriter, r *http.Request, userID uint64, vacationID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept application/json
// @Success 200 {object} api.GetVacationResponse
// @Router  /users/{user_id}/vacations/{vacation_id} [get]
func (h *handler) GetVacation(w http.ResponseWriter, r *http.Request, userID uint64, vacationID uint64) {
	ctx := r.Context()

	v, err := h.userService.GetVacation(ctx, userID, vacationID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIGetVacationResponse(v)); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
	}
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
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PutVacationJSONRequestBody true ""
// @Router  /users/{user_id}/vacations/{vacation_id} [put]
func (h *handler) PutVacation(w http.ResponseWriter, r *http.Request, userID, vacationID uint64) {
	ctx := r.Context()

	var v api.PutVacationJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &v); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := v.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	err := h.userService.UpdateVacation(ctx, userID, convert.FromAPIPutVacationRequest(vacationID, v))
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}
}
