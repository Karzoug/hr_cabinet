package middleware

import (
	"fmt"
	"net/http"

	"github.com/Employee-s-file-cabinet/backend/internal/server/errors"
)

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				errors.ReportError(r, fmt.Errorf("%s", err), true)
				errors.ErrorMessage(w, r,
					http.StatusInternalServerError,
					errors.ErrInternalServerError.Error(),
					nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
