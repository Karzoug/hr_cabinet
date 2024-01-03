package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	serr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/errors"
	uservice "github.com/Employee-s-file-cabinet/backend/internal/service/user"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

// @Produce  image/png
// @Produce  image/jpeg
// @Router   /users/{user_id}/photo [get]
func (h *handler) DownloadPhoto(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	f, closeFn, err := h.userService.DownloadPhoto(ctx, userID, r.Header.Get("If-None-Match"))
	if err != nil {
		switch {
		case errors.Is(err, uservice.ErrUserNotFound):
			serr.ErrorMessage(w, r, http.StatusNotFound, uservice.ErrUserNotFound.Error(), nil)
		case errors.Is(err, uservice.ErrPhotoFileNotFound):
			serr.ErrorMessage(w, r, http.StatusNotFound, uservice.ErrPhotoFileNotFound.Error(), nil)
		case errors.Is(err, uservice.ErrPhotoFileNotModified):
			w.WriteHeader(http.StatusNotModified)
		default:
			serr.ReportError(r, err, false)
			serr.ErrorMessage(w, r,
				http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError),
				nil)
		}
		return
	}
	defer closeFn()

	w.Header().Set("Content-Type", f.ContentType)
	if f.Hash != "" {
		w.Header().Set("ETag", f.Hash)
	}
	if _, err := io.Copy(w, f); err != nil {
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}
}

// @Accept  image/png
// @Accept  image/jpeg
// @Router  /users/{user_id}/photo [post]
func (h *handler) UploadPhoto(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	var (
		length                   int64
		isBadContentLengthHeader bool
	)
	if lengthString := r.Header.Get("Content-Length"); lengthString == "" {
		isBadContentLengthHeader = true
	} else {
		var err error
		length, err = strconv.ParseInt(lengthString, 10, 64)
		if err != nil {
			isBadContentLengthHeader = true
		}
	}
	if isBadContentLengthHeader {
		serr.ErrorMessage(w, r,
			http.StatusBadRequest,
			serr.ErrBadContentLengthHeader.Error(),
			nil)
		return
	}

	if length > uservice.MaxPhotoSize {
		serr.ErrorMessage(w, r,
			http.StatusBadRequest,
			serr.ErrLimitRequestBodySize.Error(),
			nil)
		return
	}

	lr := http.MaxBytesReader(w, r.Body, uservice.MaxPhotoSize)
	defer lr.Close()

	if err := h.userService.UploadPhoto(ctx, userID, model.File{
		Reader:      lr,
		Size:        length,
		ContentType: r.Header.Get("Content-Type"),
	}); err != nil {
		if errors.Is(err, new(http.MaxBytesError)) {
			serr.ErrorMessage(w, r,
				http.StatusBadRequest,
				serr.ErrLimitRequestBodySize.Error(),
				nil)
			return
		}
		serr.ReportError(r, err, false)
		serr.ErrorMessage(w, r,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil)
		return
	}
}
