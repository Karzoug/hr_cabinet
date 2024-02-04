package education

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

type educationHandlers struct {
	usecase user.EducationUseCase
	logger  *slog.Logger
}

func NewHandlers(e user.EducationUseCase, l *slog.Logger) educationHandlers {
	return educationHandlers{
		usecase: e,
		logger:  l,
	}
}

// @Produce application/json
// @Success 200 {object} api.ListEducationsJSONRequestBody
// @Router  /users/{user_id}/educations [get]
func (h educationHandlers) ListEducations(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	eds, err := h.usecase.List(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, toAPIListEducations(eds)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

// @Accept  application/json
// @Param   body body api.AddEducationJSONRequestBody true ""
// @Router  /users/{user_id}/educations [post]
func (h educationHandlers) AddEducation(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var e api.AddEducationJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &e); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := e.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	id, err := h.usecase.Add(ctx, userID, fromAPIAddEducationRequest(e))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+
			"/educations/"+strconv.FormatUint(id, 10))
	w.WriteHeader(http.StatusCreated)
}

// @Router /users/{user_id}/educations/{education_id} [delete]
func (h educationHandlers) DeleteEducation(w http.ResponseWriter, r *http.Request, userID uint64, educationID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.GetEducationJSONRequestBody
// @Router  /users/{user_id}/educations/{education_id} [get]
func (h educationHandlers) GetEducation(w http.ResponseWriter, r *http.Request, userID uint64, educationID uint64) {
	ctx := r.Context()

	ed, err := h.usecase.Get(ctx, userID, educationID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, toAPIGetEducationResponse(ed)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

// @Accept  application/json
// @Param   body body api.PatchEducationJSONRequestBody true ""
// @Router  /users/{user_id}/educations/{education_id} [patch]
func (h educationHandlers) PatchEducation(w http.ResponseWriter, r *http.Request, userID uint64, educationID uint64) {
	ctx := r.Context()

	var patch api.PatchEducationJSONRequestBody
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
// @Param   body body api.PutEducationJSONRequestBody true ""
// @Router  /users/{user_id}/educations/{education_id} [put]
func (h educationHandlers) PutEducation(w http.ResponseWriter, r *http.Request, userID, educationID uint64) {
	ctx := r.Context()

	var e api.PutEducationJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &e); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := e.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	err := h.usecase.Update(ctx, userID, fromAPIPutEducationRequest(educationID, e))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}
}
