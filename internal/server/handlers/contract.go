package handlers

import (
	"net/http"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
)

// @Produce application/json
// @Success 200 {array} api.Contract
// @Router  /users/{user_id}/contracts [get]
func (s *server) ListContracts(w http.ResponseWriter, r *http.Request, userID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.Contract true ""
// @Failure 409  {object} api.Error "contract already exists"
// @Router  /users/{user_id}/contracts [post]
func (s *server) AddContract(w http.ResponseWriter, r *http.Request, userID int) {
	ctx := r.Context()

	var c api.Contract
	// TODO: decode contract from request body

	if err := c.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = api.NewBadRequestErrorFromError(err)
		// encode error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /users/{user_id}/contracts/{contract_id} [delete]
func (s *server) DeleteContract(w http.ResponseWriter, r *http.Request, userID int, contractID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Contract
// @Router  /users/{user_id}/contracts/{contract_id} [get]
func (s *server) GetContract(w http.ResponseWriter, r *http.Request, userID int, contractID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PatchContractJSONRequestBody true ""
// @Router  /users/{user_id}/contracts/{contract_id} [patch]
func (s *server) PatchContract(w http.ResponseWriter, r *http.Request, userID int, contractID int) {
	ctx := r.Context()

	var patch api.PatchContractJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = api.NewBadRequestErrorFromError(err)
		// encode error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
