package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/cookie"
	srverrors "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/service/auth/model/token"

	"github.com/casbin/casbin/v2"
)

type TokenManager interface {
	Payload(token, sign string) (*token.Payload, error)
}

type Authorizer struct {
	TokenManager TokenManager
	Enforcer     *casbin.Enforcer
	logger       *slog.Logger
}

func (a *Authorizer) AuthorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(api.BearerAuthScopes) == nil {
			next.ServeHTTP(w, r)
			return
		}
		token, err := cookie.GetToken(r)
		if err != nil {
			srverrors.ResponseError(w, r,
				http.StatusForbidden,
				http.ErrNoCookie.Error(), a.logger)
			return
		}
		sign, err := cookie.GetSignature(r)
		if err != nil {
			srverrors.ResponseError(w, r,
				http.StatusForbidden,
				http.ErrNoCookie.Error(), a.logger)
			return
		}

		payload, err := a.TokenManager.Payload(token, sign)
		if err != nil {
			srverrors.ResponseError(w, r,
				http.StatusUnauthorized,
				"access token is missing or invalid", a.logger)
			return
		}

		user := payload.Data.UserID
		method := r.Method
		path := strings.TrimPrefix(r.URL.Path, api.BaseURL)

		result, _ := a.Enforcer.Enforce(user, path, method)
		if !result {
			srverrors.ResponseError(w, r,
				http.StatusUnauthorized,
				"user is not allowed to access", a.logger)
			return
		}
	})
}
