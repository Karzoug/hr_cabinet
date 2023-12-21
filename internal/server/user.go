package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/request"
	"github.com/Employee-s-file-cabinet/backend/internal/storage/s3"
)

const MaxPhotoSize = 20 << 20 // bytes

// @Produce application/json
// @Success 200 {object} api.ListUsersJSONResponseBody
// @Router  /users [get]
func (s *server) ListUsers(w http.ResponseWriter, r *http.Request, params api.ListUsersParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.FullUser true ""
// @Failure 409  {object} api.Error "user already exists"
// @Router  /users [post]
func (s *server) AddUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user api.FullUser
	// TODO: decode user from request body
	if err := user.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = api.NewBadRequestErrorFromError(err)
		// encode error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.GetUserJSONResponseBody
// @Router  /users/{user_id} [get]
func (s *server) GetUser(w http.ResponseWriter, r *http.Request, userId int, params api.GetUserParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  application/json
// @Param   body body api.PatchFullUserJSONRequestBody true ""
// @Router  /users/{user_id} [patch]
func (s *server) PatchUser(w http.ResponseWriter, r *http.Request, userID int) {
	ctx := r.Context()

	var patch api.PatchFullUserJSONRequestBody
	// TODO: decode patch from request body

	if err := patch.Validate(ctx, validator.Instance()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = api.NewBadRequestErrorFromError(err)
		// encode error
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  image/png
// @Accept  image/jpeg
// @Router  /users/{user_id}/photo [post]
func (s *server) UploadPhoto(w http.ResponseWriter, r *http.Request, userID int) {
	ctx := r.Context()

	if !request.CheckContentType([]string{"image/png", "image/jpeg"}, r.Header) {
		s.errorMessage(w, r,
			http.StatusBadRequest,
			ErrInvalidContentType.Error(),
			nil)
		return
	}

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
		s.errorMessage(w, r,
			http.StatusBadRequest,
			ErrBadContentLengthHeader.Error(),
			nil)
		return
	}

	if length > MaxPhotoSize {
		s.errorMessage(w, r,
			http.StatusBadRequest,
			ErrLimitRequestBodySize.Error(),
			nil)
		return
	}

	if exist, err := s.userRepository.ExistUser(ctx, userID); err != nil {
		s.reportServerError(r, err, false)
		s.errorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	} else if !exist {
		s.errorMessage(w, r, http.StatusNotFound, "user not found", nil)
		return
	}

	lr := http.MaxBytesReader(w, r.Body, MaxPhotoSize)
	defer lr.Close()

	if err := s.fileRepository.UploadFile(ctx, s3.File{
		Prefix:      strconv.Itoa(userID),
		Name:        "photo",
		Reader:      lr,
		Size:        length,
		ContentType: r.Header.Get("Content-Type"),
	}); err != nil {
		if errors.Is(err, new(http.MaxBytesError)) {
			s.errorMessage(w, r,
				http.StatusBadRequest,
				ErrLimitRequestBodySize.Error(),
				nil)
			return
		}
		s.reportServerError(r, err, false)
		s.errorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}
}
