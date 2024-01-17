package handlers

import (
	"net/http"
	"strconv"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/convert"
	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

// @Produce application/json
// @Success 200 {object} api.ListUsersJSONResponseBody
// @Router  /users [get]
func (h *handler) ListUsers(w http.ResponseWriter, r *http.Request, params api.ListUsersParams) {
	ctx := r.Context()

	if err := params.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
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
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	users, count, err := h.userService.ListShortUserInfo(ctx, pms)
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}
	ulist := convert.ToAPIListUsers(users)
	if err := response.JSON(w,
		http.StatusOK,
		api.ListUsersResponse{
			Users:       ulist,
			TotalUsers:  count,
			TotalPages:  (count + int(pms.Limit) - 1) / int(pms.Limit),
			CurrentPage: int(pms.Page),
		}); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
		return
	}
}

// @Accept  application/json
// @Param   body body api.AddUserJSONRequestBody true ""
// @Router  /users [post]
func (h *handler) AddUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var u api.AddUserJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &u); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := u.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	id, err := h.userService.Add(ctx, convert.FromAPIAddUserRequest(u))
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(id, 10))
	w.WriteHeader(http.StatusCreated)
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
		srverr.ResponseServiceError(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIGetUserResponse(u)); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
	}
}

func (h *handler) getExpandedUser(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	u, err := h.userService.GetExpanded(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, convert.ToAPIGetExpandedUserResponse(u)); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
	}
}

// @Accept  application/json
// @Param   body body api.PatchUserJSONRequestBody true ""
// @Router  /users/{user_id} [patch]
func (h *handler) PatchUser(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var patch api.PatchUserJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PutUserJSONRequestBody true ""
// @Router  /users/{user_id} [put]
func (h *handler) PutUser(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var u api.PutUserJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &u); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := u.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg)
		return
	}

	err := h.userService.Update(ctx, convert.FromAPIPutUserRequest(userID, u))
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}
}
