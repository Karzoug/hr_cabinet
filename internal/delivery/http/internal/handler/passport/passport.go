package passport

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

type passportHandlers struct {
	usecase user.PassportUseCase
	logger  *slog.Logger
}

func NewHandlers(e user.PassportUseCase, l *slog.Logger) passportHandlers {
	return passportHandlers{
		usecase: e,
		logger:  l,
	}
}

// @Produce application/json
// @Success 200 {object} api.ListPassportsResponse
// @Router  /users/{user_id}/passports [get]
func (h passportHandlers) ListPassports(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	psps, err := h.usecase.List(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, toAPIListPassports(psps)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

// @Accept  application/json
// @Param   body body api.AddPassportJSONRequestBody true ""
// @Failure 409  {object} api.Error "passport already exists"
// @Router  /users/{user_id}/passports [post]
func (h passportHandlers) AddPassport(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var p api.AddPassportJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &p); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := p.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	id, err := h.usecase.Add(ctx, userID, fromAPIAddPassportRequest(p))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+
			"/passports/"+strconv.FormatUint(id, 10))
	w.WriteHeader(http.StatusCreated)
}

// @Router /users/{user_id}/passports/{passport_id} [delete]
func (h passportHandlers) DeletePassport(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.GetPassportResponse
// @Router  /users/{user_id}/passports/{passport_id} [get]
func (h passportHandlers) GetPassport(w http.ResponseWriter, r *http.Request, userID, passportID uint64) {
	ctx := r.Context()

	p, err := h.usecase.Get(ctx, userID, passportID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, toAPIGetPassportResponse(p)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

// @Accept application/json
// @Param   body body api.PatchPassportJSONRequestBody true ""
// @Router  /users/{user_id}/passports/{passport_id} [patch]
func (h passportHandlers) PatchPassport(w http.ResponseWriter, r *http.Request, userID uint64, passportID uint64) {
	ctx := r.Context()

	var patch api.PatchPassportJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PutPassportJSONRequestBody true ""
// @Router  /users/{user_id}/passports/{passport_id} [put]
func (h passportHandlers) PutPassport(w http.ResponseWriter, r *http.Request, userID, passportID uint64) {
	ctx := r.Context()

	var p api.PutPassportJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &p); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := p.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	err := h.usecase.Update(ctx, userID, fromAPIPutPassportRequest(passportID, p))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}
}
