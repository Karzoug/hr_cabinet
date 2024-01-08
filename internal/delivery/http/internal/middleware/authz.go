package middleware

import (
	"net/http"
	"regexp"
	"strings"

	srvErrors "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/auth"
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
				srvErrors.ErrorMessage(w, r,
					http.StatusForbidden,
					http.ErrNoCookie.Error(), nil)
				return
			}

			ecabinetToken := cookie.Value

			payload, err := a.TokenManager.Payload(ecabinetToken)
			if err != nil {
				srvErrors.ErrorMessage(w, r,
					http.StatusUnauthorized,
					auth.ErrTokenIsInvalid.Error(), nil)
				return
			}

			user := payload.Data.UserID
			method := r.Method

			result, _ := a.Enforcer.Enforce(user, path, method)

			if !result {
				srvErrors.ErrorMessage(w, r,
					http.StatusUnauthorized,
					auth.ErrUserNotAllowed.Error(), nil)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
