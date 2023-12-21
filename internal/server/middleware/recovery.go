package middleware

import (
	"fmt"
	"net/http"

	srvErrors "github.com/Employee-s-file-cabinet/backend/internal/server/errors"
)

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				srvErrors.ReportError(r, fmt.Errorf("%s", err), true)
				srvErrors.ErrorMessage(w, r,
					http.StatusInternalServerError,
					srvErrors.ErrInternalServerError.Error(),
					nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
