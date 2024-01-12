package handlers

import (
	"errors"
	"net/http"

	"github.com/muonsoft/validation/validator"

	serr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/convert"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

// @Produce application/json
// @Success 200 {object} api.ListUsersJSONResponseBody
// @Router  /users [get]
func (h *handler) ListUsers(w http.ResponseWriter, r *http.Request, params api.ListUsersParams) {
	ctx := r.Context()

	if err := params.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	opts := make([]model.ListUsersParamsOption, 0)

	if params.Limit != nil {
		opts = append(opts, model.WithLimit(*params.Limit))
	}
	if params.Page != nil {
		opts = append(opts, model.WithPage(*params.Page))
	}
	if params.Query != nil {
		opts = append(opts, model.WithQuery(*params.Query))
	}
	if params.SortBy != nil {
		switch *params.SortBy {
		case api.ListUsersParamsSortByAlphabet:
			opts = append(opts, model.SortBy(model.ListUsersParamsSortByAlphabet))
		case api.ListUsersParamsSortByDepartment:
			opts = append(opts, model.SortBy(model.ListUsersParamsSortByDepartment))
		}
	}
	pms, err := model.NewListUsersParams(opts...)
	if err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	users, count, err := h.userService.List(ctx, *pms)
	if err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}
	ulist := convert.ToAPIShortUsers(users)
	if err := response.JSON(w,
		http.StatusOK,
		api.ListUsersJSONResponseBody{
			Users:       ulist,
			TotalUsers:  count,
			TotalPages:  (count + int(pms.Limit) - 1) / int(pms.Limit),
			CurrentPage: int(pms.Page),
		}); err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}
}

// @Accept  application/json
// @Param   body body api.FullUser true ""
// @Failure 409  {object} api.Error "user already exists"
// @Router  /users [post]
func (h *handler) AddUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user api.FullUser
	// TODO: decode user from request body
	if err := user.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.GetUserJSONResponseBody
// @Router  /users/{user_id} [get]
func (h *handler) GetUser(w http.ResponseWriter, r *http.Request, userID uint64, params api.GetUserParams) {
	if params.Expanded != nil && *params.Expanded {
		h.getExpandedUser(w, r, userID)
		return
	}
	h.getUser(w, r, userID)
}

func (h *handler) getUser(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	u, err := h.userService.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			serr.ErrorMessage(w, r, http.StatusNotFound, user.ErrUserNotFound.Error(), nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIFullUser(u)); err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
	}
}

func (h *handler) getExpandedUser(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	u, err := h.userService.GetExpanded(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			serr.ErrorMessage(w, r, http.StatusNotFound, user.ErrUserNotFound.Error(), nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIExpandedFullUser(u)); err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
	}
}

// @Accept  application/json
// @Param   body body api.PatchFullUserJSONRequestBody true ""
// @Router  /users/{user_id} [patch]
func (h *handler) PatchUser(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var patch api.PatchFullUserJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		serr.ErrorMessage(w, r, http.StatusBadRequest, msg, nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
