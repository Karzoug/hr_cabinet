package visa

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/convert"
	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
)

type Handler struct {
	service visaService
	logger  *slog.Logger
}

func NewHadlers(service visaService, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// @Produce application/json
// @Success 200 {object} api.ListVisasResponse
// @Router  /users/{user_id}/passports/{passport_id}/visas [get]
func (h *Handler) ListVisas(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	vs, err := h.service.List(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIListVisas(vs)); err != nil {
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
func (h *Handler) AddVisa(w http.ResponseWriter, r *http.Request, userID uint64) {
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

	id, err := h.service.Add(ctx, userID, convert.FromAPIAddVisaRequest(v))
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
func (h *Handler) DeleteVisa(w http.ResponseWriter, r *http.Request, userID, visaID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.GetVisaResponse
// @Router  /users/{user_id}/visas/{visa_id} [get]
func (h *Handler) GetVisa(w http.ResponseWriter, r *http.Request, userID, visaID uint64) {
	ctx := r.Context()

	v, err := h.service.Get(ctx, userID, visaID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIGetVisaResponse(v)); err != nil {
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
func (h *Handler) PatchVisa(w http.ResponseWriter, r *http.Request, userID, visaID uint64) {
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
func (h *Handler) PutVisa(w http.ResponseWriter, r *http.Request, userID, visaID uint64) {
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

	err := h.service.Update(ctx, userID, convert.FromAPIPutVisaRequest(visaID, v))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}
}
