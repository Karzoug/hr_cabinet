package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/photo"
)

const (
	errLimitRequestBodySizeMsg   = "request body too large"
	errBadContentLengthHeaderMsg = "bad content length header: missing or not a number"
)

// @Produce  image/png
// @Produce  image/jpeg
// @Router   /users/{user_id}/photo [get]
func (h *handler) DownloadPhoto(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	f, closeFn, err := h.userService.DownloadPhoto(ctx, userID, r.Header.Get("If-None-Match"))
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}
	defer closeFn()

	w.Header().Set("Content-Type", f.ContentType)
	if f.Hash != "" {
		w.Header().Set("ETag", f.Hash)
	}
	if _, err := io.Copy(w, f); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
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
		srverr.ResponseError(w, r,
			http.StatusBadRequest,
			errBadContentLengthHeaderMsg)
		return
	}

	if length > photo.MaxPhotoSize {
		srverr.ResponseError(w, r,
			http.StatusBadRequest,
			errLimitRequestBodySizeMsg)
		return
	}

	lr := http.MaxBytesReader(w, r.Body, photo.MaxPhotoSize)
	defer lr.Close()

	if err := h.userService.UploadPhoto(ctx, userID, model.File{
		Reader:      lr,
		Size:        length,
		ContentType: r.Header.Get("Content-Type"),
	}); err != nil {
		if errors.Is(err, new(http.MaxBytesError)) {
			srverr.ResponseError(w, r,
				http.StatusBadRequest,
				errLimitRequestBodySizeMsg)
			return
		}
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
		return
	}
}
