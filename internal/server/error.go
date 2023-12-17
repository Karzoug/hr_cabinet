package server

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/response"
)

var (
	ErrInvalidContentType     = errors.New("invalid content type")
	ErrLimitRequestBodySize   = errors.New("request body too large")
	ErrBadContentLengthHeader = errors.New("bad content length header: missing or not a number")
	ErrInternalServerError    = errors.New("the server encountered a problem and could not process your request")
	ErrNotFoundRoute          = errors.New("the requested resource could not be found")
	ErrMethodNotAllowed       = errors.New("the method is not supported for this resource")
)

func (s *server) reportServerError(r *http.Request, err error, withStack bool) {
	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
		trace   = string(debug.Stack())
	)

	requestAttrs := slog.Group("request", "method", method, "url", url)
	if withStack {
		s.logger.Error(message, requestAttrs, "trace", trace)
	} else {
		s.logger.Error(message, requestAttrs)
	}
}

func (s *server) errorMessage(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	message = strings.ToUpper(message[:1]) + message[1:]

	err := response.JSONWithHeaders(w, status, api.Error{Message: message}, headers)
	if err != nil {
		s.reportServerError(r, err, false)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *server) notFound(w http.ResponseWriter, r *http.Request) {
	s.errorMessage(w, r, http.StatusNotFound, ErrNotFoundRoute.Error(), nil)
}

func (s *server) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	s.errorMessage(w, r, http.StatusMethodNotAllowed, ErrMethodNotAllowed.Error(), nil)
}
