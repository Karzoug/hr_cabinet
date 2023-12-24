package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/model"
	serr "github.com/Employee-s-file-cabinet/backend/internal/server/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/request"
)

// @Accept  application/json
// @Produce application/json
// @Param   body body api.Auth true ""
// @Success 200 {object} api.Token
// @Router  /login [post]
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {git add /
	ctx := r.Context()

	var auth api.Auth
	err := request.DecodeJSON(w, r, &auth)
	if err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := auth.Validate(ctx, validator.Instance()); err != nil {
		var _ api.BadRequestError
		w.WriteHeader(http.StatusBadRequest)
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	authnData, err := h.dbRepository.GetAuthnData(ctx, auth.Login)
	if errors.Is(err, sql.ErrNoRows) {
		serr.ErrorMessage(w, r, http.StatusForbidden, serr.ErrLoginFailure.Error(), nil)
		return
	}
	if err != nil {
		serr.ErrorMessage(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	err = h.passwordVerification.Check(auth.Password, authnData.PasswordHash)
	if err != nil {
		serr.ErrorMessage(w, r, http.StatusForbidden, serr.ErrLoginFailure.Error(), nil)
		return
	}

	token, _ := h.tokenManager.Create(
		model.TokenData{
			UserID: authnData.UserID,
			RoleID: authnData.RoleID,
		})

	cookie := &http.Cookie{
		Name:     "ecabinet-token",
		Value:    token,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

// @Router /login/change-password [get]
func (h *handler) CheckKey(w http.ResponseWriter, r *http.Request, params api.CheckKeyParams) {
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
func (h *handler) InitChangePassword(w http.ResponseWriter, r *http.Request) {
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
func (h *handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
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
