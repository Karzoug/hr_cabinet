package server

import (
	"net/http"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
)

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
	w.WriteHeader(http.StatusNotImplemented)
}
