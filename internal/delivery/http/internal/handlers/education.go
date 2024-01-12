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
// @Success 200 {object} api.ListEducationsJSONRequestBody
// @Router  /users/{user_id}/educations [get]
func (h *handler) ListEducations(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	eds, err := h.userService.ListEducations(ctx, userID)
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

	if err := response.JSON(w, http.StatusOK, convert.ToAPIListEducations(eds)); err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
	}
}

// @Accept  application/json
// @Param   body body api.AddEducationJSONRequestBody true ""
// @Failure 409  {object} api.Error "education already exists"
// @Router  /users/{user_id}/educations [post]
func (h *handler) AddEducation(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var e api.AddEducationJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &e); err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := e.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	id, err := h.userService.AddEducation(ctx, userID, convert.FromAPIAddEducationRequest(e))
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

	w.Header().Set("Location", api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+"/educations/"+strconv.FormatUint(id, 10))
}

// @Router /users/{user_id}/educations/{education_id} [delete]
func (h *handler) DeleteEducation(w http.ResponseWriter, r *http.Request, userID uint64, educationID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.GetEducationJSONRequestBody
// @Router  /users/{user_id}/educations/{education_id} [get]
func (h *handler) GetEducation(w http.ResponseWriter, r *http.Request, userID uint64, educationID uint64) {
	ctx := r.Context()

	ed, err := h.userService.GetEducation(ctx, userID, educationID)
	if err != nil {
		if errors.Is(err, user.ErrEducationNotFound) {
			serr.ErrorMessage(w, r, http.StatusNotFound, user.ErrEducationNotFound.Error(), nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIGetEducationResponse(ed)); err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
	}
}

// @Accept  application/json
// @Param   body body api.PatchEducationJSONRequestBody true ""
// @Router  /users/{user_id}/educations/{education_id} [patch]
func (h *handler) PatchEducation(w http.ResponseWriter, r *http.Request, userID uint64, educationID uint64) {
	ctx := r.Context()

	var patch api.PatchEducationJSONRequestBody
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

// @Accept  application/json
// @Param   body body api.PutEducationJSONRequestBody true ""
// @Router  /users/{user_id}/educations/{education_id} [put]
func (h *handler) PutEducation(w http.ResponseWriter, r *http.Request, userID, educationID uint64) {
	ctx := r.Context()

	var e api.PutEducationJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &e); err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := e.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	err := h.userService.UpdateEducation(ctx, userID, convert.FromAPIPutEducationRequest(educationID, e))
	if err != nil {
		if errors.Is(err, user.ErrEducationNotFound) {
			serr.ErrorMessage(w, r, http.StatusNotFound, user.ErrEducationNotFound.Error(), nil)
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
