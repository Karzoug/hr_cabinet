package middleware

import (
	"fmt"
	"net/http"

	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
)

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				srverr.LogError(r, fmt.Errorf("%s", err), true)
				srverr.ResponseError(w, r,
					http.StatusInternalServerError,
					srverr.ErrInternalServerErrorMsg)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
