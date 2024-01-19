package errors

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
	"github.com/Employee-s-file-cabinet/backend/internal/service"
)

const ErrInternalServerErrorMsg = "the server encountered a problem and could not process your request"

// LogError logs the server error, with or without stack trace.
func LogError(r *http.Request, err error, withStack bool) {
	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
		trace   = string(debug.Stack())
	)

	requestAttrs := slog.Group("request", "method", method, "url", url)
	if withStack {
		slog.Error(message, requestAttrs, "trace", trace)
	} else {
		slog.Error(message, requestAttrs)
	}
}

// ResponseError converts error message to api.Error and writes this one in JSON format to response writer.
func ResponseError(w http.ResponseWriter, r *http.Request, status int, errMessage string) {
	message := strings.ToUpper(errMessage[:1]) + errMessage[1:]
	if status == http.StatusNotModified {
		// RFC 2616:
		// The 304 response MUST NOT contain a message-body,
		// and thus is always terminated by the first empty line after the header fields.
		w.WriteHeader(status)
		return
	}
	if err := response.JSON(w,
		status,
		api.Error{Message: message}); err != nil {
		LogError(r, err, false)
	}
}

// ResponseServiceError converts a service error to api.Error and writes this one in JSON format to response writer.
func ResponseServiceError(w http.ResponseWriter, r *http.Request, err error) {
	serviceErr := new(service.Error)
	if !errors.As(err, &serviceErr) {
		LogError(r, err, false)
		ResponseError(w, r,
			http.StatusInternalServerError,
			ErrInternalServerErrorMsg)
		return
	}
	ResponseError(w, r,
		serviceStatusToHTTPStatusCode(serviceErr),
		serviceErr.Error())
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	ResponseError(w, r, http.StatusNotFound, "the requested resource could not be found")
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	ResponseError(w, r, http.StatusMethodNotAllowed, "the method is not supported for this resource")
}

func serviceStatusToHTTPStatusCode(err *service.Error) int {
	switch err.Status {
	case service.NotFound:
		return http.StatusNotFound
	case service.InvalidArgument:
		return http.StatusBadRequest
	case service.AlreadyExists:
		return http.StatusBadRequest
	case service.NotModified:
		return http.StatusNotModified
	case service.Conflict:
		return http.StatusConflict
	case service.PermissionDenied:
		return http.StatusForbidden
	case service.Unauthenticated:
		return http.StatusUnauthorized
	case service.ContentTooLarge:
		return http.StatusRequestEntityTooLarge
	default:
		return http.StatusInternalServerError
	}
}
