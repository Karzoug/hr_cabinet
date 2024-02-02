package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
)

func RecoverPanicFn(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					srverr.LogError(r, fmt.Errorf("%s", err), true, logger)
					srverr.ResponseError(w, r,
						http.StatusInternalServerError,
						srverr.ErrInternalServerErrorMsg, logger)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
