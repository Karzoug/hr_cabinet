package contract

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

type contractHandlers struct {
	usecase user.ContractUseCase
	logger  *slog.Logger
}

func NewHandlers(c user.ContractUseCase, l *slog.Logger) contractHandlers {
	return contractHandlers{
		usecase: c,
		logger:  l,
	}
}

// @Produce application/json
// @Success 200 {object} api.ListContractsJSONRequestBody
// @Router  /users/{user_id}/contracts [get]
func (h contractHandlers) ListContracts(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	cs, err := h.usecase.List(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, toAPIListContracts(cs)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

// @Accept  application/json
// @Param   body body api.AddContractJSONRequestBody true ""
// @Failure 409  {object} api.Error "contract already exists"
// @Router  /users/{user_id}/contracts [post]
func (h contractHandlers) AddContract(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var ct api.AddContractJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &ct); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := ct.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	id, err := h.usecase.Add(ctx, userID, fromAPIAddContractRequest(ct))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+
			"/contracts/"+strconv.FormatUint(id, 10))
}

// @Router /users/{user_id}/contracts/{contract_id} [delete]
func (h contractHandlers) DeleteContract(w http.ResponseWriter, r *http.Request, userID uint64, contractID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.GetContractJSONRequestBody
// @Router  /users/{user_id}/contracts/{contract_id} [get]
func (h contractHandlers) GetContract(w http.ResponseWriter, r *http.Request, userID uint64, contractID uint64) {
	ctx := r.Context()

	c, err := h.usecase.Get(ctx, userID, contractID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, toAPIGetContractResponse(c)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

// @Accept  application/json
// @Param   body body api.PatchContractJSONRequestBody true ""
// @Router  /users/{user_id}/contracts/{contract_id} [patch]
func (h contractHandlers) PatchContract(w http.ResponseWriter, r *http.Request, userID uint64, contractID uint64) {
	ctx := r.Context()

	var patch api.PatchContractJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PutContractJSONRequestBody true ""
// @Router  /users/{user_id}/contracts/{contract_id} [put]
func (h contractHandlers) PutContract(w http.ResponseWriter, r *http.Request, userID uint64, contractID uint64) {
	ctx := r.Context()

	var c api.PutContractJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &c); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := c.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	err := h.usecase.Update(ctx, userID, fromAPIPutContractRequest(contractID, c))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}
}
