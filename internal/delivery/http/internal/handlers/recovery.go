package handlers

import (
	"errors"
	"net/http"

	"github.com/muonsoft/validation/validator"

	srvErrors "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

// @Accept  application/json
// @Param   body body api.InitChangePasswordRequest true ""
// @Router  /login/init-change-password [post]
func (h *handler) InitChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var initChngPswdReq api.InitChangePasswordRequest

	err := request.DecodeJSON(w, r, &initChngPswdReq)
	if err != nil {
		srvErrors.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := initChngPswdReq.Validate(ctx, validator.Instance()); err != nil {
		srvErrors.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err = h.passwordRecoveryService.InitChangePassword(ctx, string(initChngPswdReq.Login)); err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotFound):
			srvErrors.ErrorMessage(w, r,
				http.StatusNotFound,
				http.StatusText(http.StatusNotFound),
				nil)
		default:
			srvErrors.ReportError(r, err, false)
			srvErrors.ErrorMessage(w, r,
				http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				nil)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Accept  application/json
// @Param   body body api.ChangePasswordRequest true ""
// @Router  /login/change-password [post]
func (h *handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var chPsw api.ChangePasswordRequest

	err := request.DecodeJSON(w, r, &chPsw)
	if err != nil {
		srvErrors.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := chPsw.Validate(ctx, validator.Instance()); err != nil {
		srvErrors.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = h.passwordRecoveryService.ChangePassword(ctx, chPsw.Key, chPsw.Password)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotFound):
			srvErrors.ErrorMessage(w, r,
				http.StatusNotFound,
				http.StatusText(http.StatusNotFound),
				nil)
		default:
			srvErrors.ReportError(r, err, false)
			srvErrors.ErrorMessage(w, r,
				http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				nil)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Router /login/change-password [get]
func (h *handler) CheckKey(w http.ResponseWriter, r *http.Request, params api.CheckKeyParams) {
	// TODO: ограничение количества запросов
	ctx := r.Context()

	if err := params.Validate(ctx, validator.Instance()); err != nil {
		srvErrors.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.passwordRecoveryService.Check(ctx, params.Key); err != nil {
		srvErrors.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusOK)
}
