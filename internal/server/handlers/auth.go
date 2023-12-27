package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/model"
	serr "github.com/Employee-s-file-cabinet/backend/internal/server/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/request"
)

// TODO: перенести в переменые окружения
const domen = "https://ecabinet.acceleratorpracticum.ru"

// @Accept  application/json
// @Produce application/json
// @Param   body body api.Auth true ""
// @Success 200 {object} api.Token
// @Router  /login [post]
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var auth api.Auth
	err := request.DecodeJSON(w, r, &auth)
	if err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := auth.Validate(ctx, validator.Instance()); err != nil {
		var _ api.BadRequestError
		w.WriteHeader(http.StatusBadRequest)
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	authnData, err := h.dbRepository.GetAuthnData(ctx, auth.Login)
	if errors.Is(err, sql.ErrNoRows) {
		serr.ErrorMessage(w, r, http.StatusForbidden, serr.ErrLoginFailure.Error(), nil)
		return
	}
	if err != nil {
		serr.ErrorMessage(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	err = h.passwordVerification.Check(auth.Password, authnData.PasswordHash)
	if err != nil {
		serr.ErrorMessage(w, r, http.StatusForbidden, serr.ErrLoginFailure.Error(), nil)
		return
	}

	token, _ := h.tokenManager.Create(
		model.TokenData{
			UserID: authnData.UserID,
			RoleID: authnData.RoleID,
		})

	cookie := &http.Cookie{
		Name:     "ecabinet-token",
		Value:    token,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

// @Router /login/change-password [get]
func (h *handler) CheckKey(w http.ResponseWriter, r *http.Request, params api.CheckKeyParams) {
	// TODO: ограничение количества запросов
	ctx := r.Context()

	if err := params.Validate(ctx, validator.Instance()); err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	//проверка наличия и срока действия ключа
	login, err := h.keyRepository.Get(ctx, params.Key)
	if err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if login == "" {
		serr.ErrorMessage(w, r, http.StatusInternalServerError, "login lost", nil)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Accept  application/json
// @Param   body body api.InitChangePasswordRequest true ""
// @Router  /login/init-change-password [post]
func (h *handler) InitChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var chPsw api.InitChangePasswordRequest

	err := request.DecodeJSON(w, r, &chPsw)
	if err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := chPsw.Validate(ctx, validator.Instance()); err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	//обращение к базе, проверка наличия пользователя с заданным логином
	exist, err := h.dbRepository.ExistEmployee(ctx, chPsw.Login)
	if err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}
	if !exist {
		serr.ErrorMessage(w, r, http.StatusNotFound, "employee not found", nil)
		return
	}

	//генерация ключа
	randBytes := make([]byte, 26)
	_, err = rand.Read(randBytes)
	if err != nil {
		serr.ErrorMessage(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	randString := base64.StdEncoding.EncodeToString(randBytes)

	//сохранение ключа в мапе
	err = h.keyRepository.Set(ctx, randString, chPsw.Login, time.Minute*30)
	if err != nil {
		serr.ErrorMessage(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	//отправка письма
	subject := "Запрос на восстановление доступа"
	msg := fmt.Sprintf(`Для восстановления доступа к личному кабинету перейдите по ссылке:
	%s/access-restore/password-reset?key=%s`, domen, randString)

	if err := h.mail.SendSSLMail(subject, msg, chPsw.Login); err != nil {
		serr.ErrorMessage(w, r, http.StatusInternalServerError, err.Error(), nil)
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
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := chPsw.Validate(ctx, validator.Instance()); err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	//получение пользователя по ключу
	login, err := h.keyRepository.Get(ctx, chPsw.Key)
	if err != nil {
		serr.ErrorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if login == "" {
		serr.ErrorMessage(w, r, http.StatusInternalServerError, "login lost", nil)
		return
	}

	//TODO: проверка пароля на сложность

	passHash, err := h.passwordVerification.Hash(chPsw.Password)
	if err != nil {
		serr.ErrorMessage(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	err = h.dbRepository.ChangePass(ctx, login, passHash)
	if err != nil {
		//TODO: анализировать виды ошибок
		serr.ErrorMessage(w, r, http.StatusNotFound, "employee not found", nil)
		return
	}

	w.WriteHeader(http.StatusOK)
}
