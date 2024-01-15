package handlers

import (
	"errors"
	"net/http"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
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
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := initChngPswdReq.Validate(ctx, validator.Instance()); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.passwordRecoveryService.InitChangePassword(ctx, string(initChngPswdReq.Login)); err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotFound):
			srverr.ResponseError(w, r,
				http.StatusNotFound,
				http.StatusText(http.StatusNotFound))
		default:
			srverr.LogError(r, err, false)
			srverr.ResponseError(w, r,
				http.StatusInternalServerError,
				srverr.ErrInternalServerErrorMsg)
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
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := chPsw.Validate(ctx, validator.Instance()); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = h.passwordRecoveryService.ChangePassword(ctx, chPsw.Key, chPsw.Password)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotFound):
			srverr.ResponseError(w, r,
				http.StatusNotFound,
				http.StatusText(http.StatusNotFound))
		default:
			srverr.LogError(r, err, false)
			srverr.ResponseError(w, r,
				http.StatusInternalServerError,
				srverr.ErrInternalServerErrorMsg)
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
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.passwordRecoveryService.Check(ctx, params.Key); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
