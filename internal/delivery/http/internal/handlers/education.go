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
// @Success 200 {object} api.ListEducationsJSONRequestBody
// @Router  /users/{user_id}/educations [get]
func (h *handler) ListEducations(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	eds, err := h.userService.ListEducations(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIListEducations(eds)); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
	}
}

// @Accept  application/json
// @Param   body body api.AddEducationJSONRequestBody true ""
// @Router  /users/{user_id}/educations [post]
func (h *handler) AddEducation(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var e api.AddEducationJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &e); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := e.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	id, err := h.userService.AddEducation(ctx, userID, convert.FromAPIAddEducationRequest(e))
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+
			"/educations/"+strconv.FormatUint(id, 10))
	w.WriteHeader(http.StatusCreated)
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
		srverr.ResponseServiceError(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIGetEducationResponse(ed)); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
	}
}

// @Accept  application/json
// @Param   body body api.PatchEducationJSONRequestBody true ""
// @Router  /users/{user_id}/educations/{education_id} [patch]
func (h *handler) PatchEducation(w http.ResponseWriter, r *http.Request, userID uint64, educationID uint64) {
	ctx := r.Context()

	var patch api.PatchEducationJSONRequestBody
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
// @Param   body body api.PutEducationJSONRequestBody true ""
// @Router  /users/{user_id}/educations/{education_id} [put]
func (h *handler) PutEducation(w http.ResponseWriter, r *http.Request, userID, educationID uint64) {
	ctx := r.Context()

	var e api.PutEducationJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &e); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := e.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	err := h.userService.UpdateEducation(ctx, userID, convert.FromAPIPutEducationRequest(educationID, e))
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}
}
