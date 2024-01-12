package handlers

import (
	"errors"
	"net/http"

	"github.com/muonsoft/validation/validator"

	srvErrors "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
	authErrors "github.com/Employee-s-file-cabinet/backend/internal/service/auth"
)

// @Accept  application/json
// @Produce application/json
// @Param   body body api.LoginJSONRequestBody true ""
// @Router  /login [post]
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var auth api.LoginJSONRequestBody
	err := request.DecodeJSON(w, r, &auth)
	if err != nil {
		srvErrors.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := auth.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		srvErrors.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token, err := h.authService.Login(ctx, string(auth.Login), auth.Password)
	if err != nil {
		switch {
		case errors.Is(err, authErrors.ErrForbidden):
			srvErrors.ErrorMessage(w, r,
				http.StatusForbidden,
				srvErrors.ErrLoginFailure.Error(), nil)
		default:
			srvErrors.ReportError(r, err, false)
			srvErrors.ErrorMessage(w, r,
				http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), nil)
		}
		return
	}

	cookie := &http.Cookie{
		Name:     "ecabinet-token",
		Value:    token,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Expires:  h.authService.Expires(),
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}
