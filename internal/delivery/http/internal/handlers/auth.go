package handlers

import (
	"net/http"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
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
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := auth.Validate(ctx, validator.Instance()); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authService.Login(ctx, string(auth.Login), auth.Password)
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
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
