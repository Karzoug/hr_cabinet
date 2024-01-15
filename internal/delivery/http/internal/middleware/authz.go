package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	srverrors "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/service/auth/model/token"

	"github.com/casbin/casbin/v2"
)

const (
	cookieName = "ecabinet-token"
	pattern    = `^/login|^/health`
)

type TokenManager interface {
	Payload(token string) (*token.Payload, error)
}

type Authorizer struct {
	TokenManager TokenManager
	Enforcer     *casbin.Enforcer
}

func (a *Authorizer) AuthorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, api.BaseURL)

		matched, _ := regexp.MatchString(pattern, path)

		if !matched {
			cookie, err := r.Cookie(cookieName)
			if err != nil {
				srverrors.ResponseError(w, r,
					http.StatusForbidden,
					http.ErrNoCookie.Error())
				return
			}

			ecabinetToken := cookie.Value

			payload, err := a.TokenManager.Payload(ecabinetToken)
			if err != nil {
				srverrors.ResponseError(w, r,
					http.StatusUnauthorized,
					"access token is missing or invalid")
				return
			}

			user := payload.Data.UserID
			method := r.Method

			result, _ := a.Enforcer.Enforce(user, path, method)

			if !result {
				srverrors.ResponseError(w, r,
					http.StatusUnauthorized,
					"user is not allowed to access")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
