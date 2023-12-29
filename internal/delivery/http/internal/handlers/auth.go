package handlers

import (
	"errors"
	"net/http"

	"github.com/muonsoft/validation/validator"

	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
	autherr "github.com/Employee-s-file-cabinet/backend/internal/service/auth"
)

// @Accept  application/json
// @Produce application/json
// @Param   body body api.Auth true ""
// @Success 200 {object} api.Token
// @Router  /login [post]
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var auth api.Auth
	err := request.DecodeJSON(w, r, &auth)
	if err != nil {
		srverr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := auth.Validate(ctx, validator.Instance()); err != nil {
		var _ api.BadRequestError
		w.WriteHeader(http.StatusBadRequest)
		srverr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token, err := h.authService.Login(ctx, auth.Login, auth.Password)
	if err != nil {
		switch {
		case errors.Is(err, autherr.ErrForbidden):
			srverr.ErrorMessage(w, r,
				http.StatusForbidden,
				srverr.ErrLoginFailure.Error(), nil)
		default:
			srverr.ReportError(r, err, false)
			srverr.ErrorMessage(w, r,
				http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), nil)
		}
	}

	cookie := &http.Cookie{
		Name:     "ecabinet-token",
		Value:    token,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		Expires:  h.authService.Expires(),
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
