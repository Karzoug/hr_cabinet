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
// @Success 200 {object} api.ListContractsJSONRequestBody
// @Router  /users/{user_id}/contracts [get]
func (h *handler) ListContracts(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	cs, err := h.userService.ListContracts(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIListContracts(cs)); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
	}
}

// @Accept  application/json
// @Param   body body api.AddContractJSONRequestBody true ""
// @Failure 409  {object} api.Error "contract already exists"
// @Router  /users/{user_id}/contracts [post]
func (h *handler) AddContract(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var ct api.AddContractJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &ct); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := ct.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	id, err := h.userService.AddContract(ctx, userID, convert.FromAPIAddContractRequest(ct))
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+
			"/contracts/"+strconv.FormatUint(id, 10))
}

// @Router /users/{user_id}/contracts/{contract_id} [delete]
func (h *handler) DeleteContract(w http.ResponseWriter, r *http.Request, userID uint64, contractID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.GetContractJSONRequestBody
// @Router  /users/{user_id}/contracts/{contract_id} [get]
func (h *handler) GetContract(w http.ResponseWriter, r *http.Request, userID uint64, contractID uint64) {
	ctx := r.Context()

	c, err := h.userService.GetContract(ctx, userID, contractID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIGetContractResponse(c)); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
	}
}

// @Accept  application/json
// @Param   body body api.PatchContractJSONRequestBody true ""
// @Router  /users/{user_id}/contracts/{contract_id} [patch]
func (h *handler) PatchContract(w http.ResponseWriter, r *http.Request, userID uint64, contractID uint64) {
	ctx := r.Context()

	var patch api.PatchContractJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PutContractJSONRequestBody true ""
// @Router  /users/{user_id}/contracts/{contract_id} [put]
func (h *handler) PutContract(w http.ResponseWriter, r *http.Request, userID uint64, contractID uint64) {
	ctx := r.Context()

	var c api.PutContractJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &c); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	err := h.userService.UpdateContract(ctx, userID, convert.FromAPIPutContractRequest(contractID, c))
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}
}
