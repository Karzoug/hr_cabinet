package handlers

import (
	"net/http"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
)

// @Accept  application/json
// @Produce application/json
// @Param   body body api.Auth true ""
// @Success 200 {object} api.Token
// @Router  /login [post]
func (s *server) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var auth api.Auth
	// TODO: decode auth from request body

	if err := auth.Validate(ctx, validator.Instance()); err != nil {
		var _ api.BadRequestError
		w.WriteHeader(http.StatusBadRequest)
		// TODO: return error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /login/change-password [get]
func (s *server) CheckKey(w http.ResponseWriter, r *http.Request, params api.CheckKeyParams) {
	ctx := r.Context()

	if err := params.Validate(ctx, validator.Instance()); err != nil {
		var _ api.BadRequestError
		w.WriteHeader(http.StatusBadRequest)
		// TODO: return error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.InitChangePasswordRequest true ""
// @Router  /login/init-change-password [post]
func (s *server) InitChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var chPsw api.InitChangePasswordRequest

	// TODO: decode InitChangePasswordRequest from request body

	if err := chPsw.Validate(ctx, validator.Instance()); err != nil {
		var _ api.BadRequestError
		w.WriteHeader(http.StatusBadRequest)
		// TODO: return error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.ChangePasswordRequest true ""
// @Router  /login/change-password [post]
func (s *server) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var chPsw api.ChangePasswordRequest

	// TODO: decode ChangePasswordRequest from request body

	if err := chPsw.Validate(ctx, validator.Instance()); err != nil {
		var _ api.BadRequestError
		w.WriteHeader(http.StatusBadRequest)
		// TODO: return error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
