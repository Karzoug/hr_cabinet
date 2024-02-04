package visa

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user"
	"github.com/muonsoft/validation/validator"
)

type visaHandlers struct {
	usecase user.VisaUseCase
	logger  *slog.Logger
}

func NewVisaHandlers(v user.VisaUseCase, l *slog.Logger) visaHandlers {
	return visaHandlers{
		usecase: v,
		logger:  l,
	}
}

// @Produce application/json
// @Success 200 {object} api.ListVisasResponse
// @Router  /users/{user_id}/passports/{passport_id}/visas [get]
func (h visaHandlers) ListVisas(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	vs, err := h.usecase.List(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, toAPIListVisas(vs)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

// @Accept application/json
// @Param   body body api.AddVisaJSONRequestBody true ""
// @Router  /users/{user_id}/visas [post]
func (h visaHandlers) AddVisa(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var v api.AddVisaJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &v); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := v.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	id, err := h.usecase.Add(ctx, userID, fromAPIAddVisaRequest(v))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+
			"/visas/"+strconv.FormatUint(id, 10))
	w.WriteHeader(http.StatusCreated)
}

// @Router /users/{user_id}/visas/{visa_id} [delete]
func (h visaHandlers) DeleteVisa(w http.ResponseWriter, r *http.Request, userID, visaID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.GetVisaResponse
// @Router  /users/{user_id}/visas/{visa_id} [get]
func (h visaHandlers) GetVisa(w http.ResponseWriter, r *http.Request, userID, visaID uint64) {
	ctx := r.Context()

	v, err := h.usecase.Get(ctx, userID, visaID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, toAPIGetVisaResponse(v)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

// @Accept application/json
// @Param   body body api.PatchVisaJSONRequestBody true ""
// @Router  /users/{user_id}/visas/{visa_id} [patch]
func (h visaHandlers) PatchVisa(w http.ResponseWriter, r *http.Request, userID, visaID uint64) {
	ctx := r.Context()

	var patch api.PatchVisaJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PutVisaJSONRequestBody true ""
// @Router  /users/{user_id}/visas/{visa_id} [put]
func (h visaHandlers) PutVisa(w http.ResponseWriter, r *http.Request, userID, visaID uint64) {
	ctx := r.Context()

	var v api.PutVisaJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &v); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := v.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	err := h.usecase.Update(ctx, userID, fromAPIPutVisaRequest(visaID, v))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}
}
