package cookie

import (
	"net/http"
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/config/env"
)

const (
	tokenName = "ecabinet-token"
	signName  = "ecabinet-token-sign"
)

func GetToken(r *http.Request) (string, error) {
	return get(r, tokenName)
}

func GetSignature(r *http.Request) (string, error) {
	return get(r, signName)
}

func SetToken(w http.ResponseWriter, token string, expires time.Time, envType env.Type) {
	cookie := &http.Cookie{
		Name:     tokenName,
		Value:    token,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		Expires:  expires,
	}
	if envType == env.Development {
		cookie.SameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, cookie)
}

func SetSignature(w http.ResponseWriter, sign string, expires time.Time, envType env.Type) {
	cookie := &http.Cookie{
		Name:     signName,
		Value:    sign,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
		Expires:  expires,
	}
	if envType == env.Development {
		cookie.SameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, cookie)
}

func get(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
