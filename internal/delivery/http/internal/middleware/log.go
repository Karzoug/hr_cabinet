package middleware

import (
	"log/slog"
	"net/http"

	"github.com/tomasen/realip"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
)

func LogAccessFn(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mw := response.NewMetricsResponseWriter(w)
			next.ServeHTTP(mw, r)

			var (
				ip     = realip.FromRequest(r)
				method = r.Method
				url    = r.URL.String()
				proto  = r.Proto
			)

			userAttrs := slog.Group("user", "ip", ip)
			requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto)
			responseAttrs := slog.Group("response", "status", mw.StatusCode, "size", mw.BytesCount)

			logger.Info("access", userAttrs, requestAttrs, responseAttrs)
		})
	}
}
