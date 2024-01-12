package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/muonsoft/validation/validator"

	serr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/convert"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user"
)

// @Produce application/json
// @Success 200 {object} api.ListVacationsResponse
// @Router  /users/{user_id}/vacations [get]
func (h *handler) ListVacations(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	vacations, err := h.userService.ListVacations(ctx, userID)
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

	if err := response.JSON(w, http.StatusOK, convert.ToAPIListVacations(vacations)); err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
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
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := v.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	id, err := h.userService.AddVacation(ctx, userID, convert.FromAPIAddVacationRequest(v))
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
			"/vacations/"+strconv.FormatUint(id, 10))
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
		if errors.Is(err, user.ErrVacationNotFound) {
			serr.ErrorMessage(w, r, http.StatusNotFound, user.ErrVacationNotFound.Error(), nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIGetVacationResponse(v)); err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
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
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
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
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := v.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	err := h.userService.UpdateVacation(ctx, userID, convert.FromAPIPutVacationRequest(vacationID, v))
	if err != nil {
		if errors.Is(err, user.ErrVacationNotFound) {
			serr.ErrorMessage(w, r, http.StatusNotFound, user.ErrVacationNotFound.Error(), nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}
}
