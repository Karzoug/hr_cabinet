package handlers

import (
	"net/http"

	"github.com/muonsoft/validation/validator"

	serr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
)

// @Produce application/json
// @Success 200 {object} api.ListContractsJSONRequestBody
// @Router  /users/{user_id}/contracts [get]
func (h *handler) ListContracts(w http.ResponseWriter, r *http.Request, userID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.AddContractJSONRequestBody true ""
// @Failure 409  {object} api.Error "contract already exists"
// @Router  /users/{user_id}/contracts [post]
func (h *handler) AddContract(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var c api.AddContractJSONRequestBody
	// TODO: decode contract from request body

	if err := c.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /users/{user_id}/contracts/{contract_id} [delete]
func (h *handler) DeleteContract(w http.ResponseWriter, r *http.Request, userID uint64, contractID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.GetContractJSONRequestBody
// @Router  /users/{user_id}/contracts/{contract_id} [get]
func (h *handler) GetContract(w http.ResponseWriter, r *http.Request, userID uint64, contractID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
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
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PutContractJSONRequestBody true ""
// @Router  /users/{user_id}/contracts/{contract_id} [put]
func (h *handler) PutContract(w http.ResponseWriter, r *http.Request, userID uint64, contractID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}
