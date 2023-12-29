package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/muonsoft/validation/validator"
	"github.com/oapi-codegen/runtime/types"

	serr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
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
	ulist := convertUsersToAPIShortUsers(users)
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

	if err := response.JSON(w, http.StatusOK, toAPIFullUser(u)); err != nil {
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

// @Accept  image/png
// @Accept  image/jpeg
// @Router  /users/{user_id}/photo [post]
func (h *handler) UploadPhoto(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var (
		length                   int64
		isBadContentLengthHeader bool
	)
	if lengthString := r.Header.Get("Content-Length"); lengthString == "" {
		isBadContentLengthHeader = true
	} else {
		var err error
		length, err = strconv.ParseInt(lengthString, 10, 64)
		if err != nil {
			isBadContentLengthHeader = true
		}
	}
	if isBadContentLengthHeader {
		serr.ErrorMessage(w, r,
			http.StatusBadRequest,
			serr.ErrBadContentLengthHeader.Error(),
			nil)
		return
	}

	if length > user.MaxPhotoSize {
		serr.ErrorMessage(w, r,
			http.StatusBadRequest,
			serr.ErrLimitRequestBodySize.Error(),
			nil)
		return
	}

	lr := http.MaxBytesReader(w, r.Body, user.MaxPhotoSize)
	defer lr.Close()

	if err := h.userService.UploadPhoto(ctx, userID, model.File{
		Reader:      lr,
		Size:        length,
		ContentType: r.Header.Get("Content-Type"),
	}); err != nil {
		if errors.Is(err, new(http.MaxBytesError)) {
			serr.ErrorMessage(w, r,
				http.StatusBadRequest,
				serr.ErrLimitRequestBodySize.Error(),
				nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}
}

func toAPIFullUser(u *model.User) api.FullUser {
	var gr api.Gender
	switch u.Gender {
	case model.GenderFemale:
		gr = api.GenderFemale
	case model.GenderMale:
		gr = api.GenderMale
	}
	return api.FullUser{
		ShortUser: convertUserToAPIShortUser(u),
		DateOfBirth: types.Date{
			Time: u.DateOfBirth,
		},
		Gender: gr,
		Grade:  u.Grade,
		Insurance: api.Insurance{
			Number: u.InsuranceNumber,
		},
		Nationality:         u.Nationality,
		PlaceOfBirth:        u.PlaceOfBirth,
		RegistrationAddress: u.RegistrationAddress,
		ResidentialAddress:  u.ResidentialAddress,
		Taxpayer: api.Taxpayer{
			Number: u.TaxpayerNumber,
		},
	}
}

func convertUsersToAPIShortUsers(users []model.User) []api.ShortUser {
	res := make([]api.ShortUser, len(users))
	for i := 0; i < len(users); i++ {
		res[i] = convertUserToAPIShortUser(&users[i])
	}
	return res
}

func convertUserToAPIShortUser(u *model.User) api.ShortUser {
	var phn map[string]api.PhoneNumber
	if u.PhoneNumbers != nil {
		phn = make(map[string]api.PhoneNumber, len(u.PhoneNumbers))
		for k, v := range u.PhoneNumbers {
			phn[k] = api.PhoneNumber(v)
		}
	}
	return api.ShortUser{
		ID:           &u.ID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		MiddleName:   u.MiddleName,
		Email:        u.Email,
		Position:     u.Position,
		Department:   u.Department,
		PhoneNumbers: phn,
	}
}
