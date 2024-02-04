package user

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/request"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type userHandlers struct {
	usecase user.UserUseCase
	logger  *slog.Logger
}

func newUserHandlers(t user.UserUseCase, l *slog.Logger) userHandlers {
	return userHandlers{
		usecase: t,
		logger:  l,
	}
}

// @Produce application/json
// @Success 200 {object} api.ListUsersJSONResponseBody
// @Router  /users [get]
func (h userHandlers) ListUsers(w http.ResponseWriter, r *http.Request, params api.ListUsersParams) {
	ctx := r.Context()

	if err := params.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
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
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	users, count, err := h.usecase.ListShortUserInfo(ctx, pms)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}
	ulist := toAPIListUsers(users)
	if err := response.JSON(w,
		http.StatusOK,
		api.ListUsersResponse{
			Users:       ulist,
			TotalUsers:  count,
			TotalPages:  (count + int(pms.Limit) - 1) / int(pms.Limit),
			CurrentPage: int(pms.Page),
		}); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
		return
	}
}

// @Accept  application/json
// @Param   body body api.AddUserJSONRequestBody true ""
// @Router  /users [post]
func (h userHandlers) AddUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var u api.AddUserJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &u); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := u.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	id, err := h.usecase.Add(ctx, fromAPIAddUserRequest(u))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(id, 10))
	w.WriteHeader(http.StatusCreated)
}

// @Produce application/json
// @Success 200 {object} api.GetUserJSONResponseBody
// @Router  /users/{user_id} [get]
func (h userHandlers) GetUser(w http.ResponseWriter, r *http.Request, userID uint64, params api.GetUserParams) {
	if params.Expanded != nil && *params.Expanded {
		h.getExpandedUser(w, r, userID)
		return
	}
	h.getUser(w, r, userID)
}

func (h userHandlers) getUser(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	u, err := h.usecase.Get(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, toAPIGetUserResponse(u)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

func (h userHandlers) getExpandedUser(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	u, err := h.usecase.GetExpanded(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err := response.JSON(w, http.StatusOK, toAPIGetExpandedUserResponse(u)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

// @Accept  application/json
// @Param   body body api.PatchUserJSONRequestBody true ""
// @Router  /users/{user_id} [patch]
func (h userHandlers) PatchUser(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var patch api.PatchUserJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PutUserJSONRequestBody true ""
// @Router  /users/{user_id} [put]
func (h userHandlers) PutUser(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var u api.PutUserJSONRequestBody
	if err := request.DecodeJSONStrict(w, r, &u); err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := u.Validate(ctx, validator.Instance()); err != nil {
		msg := api.ValidationErrorMessage(err)
		srverr.ResponseError(w, r, http.StatusBadRequest, msg, h.logger)
		return
	}

	err := h.usecase.Update(ctx, fromAPIPutUserRequest(userID, u))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}
}
