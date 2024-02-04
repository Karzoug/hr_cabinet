package photo

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

const (
	errLimitRequestBodySizeMsg   = "request body too large"
	errBadContentLengthHeaderMsg = "bad content length header: missing or not a number"
)

type photoHandlers struct {
	usecase user.PhotoUseCase
	logger  *slog.Logger
}

func newPhotoHandlers(t user.PhotoUseCase, l *slog.Logger) photoHandlers {
	return photoHandlers{
		usecase: t,
		logger:  l,
	}
}

// @Produce  image/png
// @Produce  image/jpeg
// @Router   /users/{user_id}/photo [get]
func (h photoHandlers) DownloadPhoto(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	f, closeFn, err := h.usecase.Download(ctx, userID, r.Header.Get("If-None-Match"))
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}
	defer closeFn()

	w.Header().Set("Content-Type", f.ContentType)
	if f.Hash != "" {
		w.Header().Set("ETag", f.Hash)
	}
	if _, err := io.Copy(w, f); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
		return
	}
}

// @Accept  image/png
// @Accept  image/jpeg
// @Router  /users/{user_id}/photo [post]
func (h photoHandlers) UploadPhoto(w http.ResponseWriter, r *http.Request, userID uint64) {
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
			errBadContentLengthHeaderMsg,
			h.logger)
		return
	}

	if length > user.MaxPhotoSize {
		srverr.ResponseError(w, r,
			http.StatusBadRequest,
			errLimitRequestBodySizeMsg,
			h.logger)
		return
	}

	lr := http.MaxBytesReader(w, r.Body, user.MaxPhotoSize)
	defer lr.Close()

	if err := h.usecase.Upload(ctx, userID, model.File{
		Reader:      lr,
		Size:        length,
		ContentType: r.Header.Get("Content-Type"),
	}); err != nil {
		if errors.Is(err, new(http.MaxBytesError)) {
			srverr.ResponseError(w, r,
				http.StatusBadRequest,
				errLimitRequestBodySizeMsg,
				h.logger)
			return
		}
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
		return
	}
}
