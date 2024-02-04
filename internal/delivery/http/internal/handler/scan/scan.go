package scan

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

const errLimitRequestBodySizeMsg = "request body too large"

type scanHandlers struct {
	usecase user.ScanUseCase
	logger  *slog.Logger
}

func NewHandlers(s user.ScanUseCase, l *slog.Logger) scanHandlers {
	return scanHandlers{
		usecase: s,
		logger:  l,
	}
}

// @Produce application/json
// @Success 200 {array} api.Scan
// @Router  /users/{user_id}/scans [get]
func (h scanHandlers) ListScans(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	scans, err := h.usecase.List(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err = response.JSON(w, http.StatusOK, toAPIScans(scans)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}

// @Accept  multipart/form-data
// @Param   body body api.UploadScanMultipartRequestBody true ""
// @Router  /users/{user_id}/scans [post]
func (h scanHandlers) UploadScan(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	sc, err := handleScanMultipartRequest(ctx, r)
	if err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	file, header, err := r.FormFile("fileName")
	if err != nil {
		srverr.ResponseError(w, r, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	if header.Size > user.MaxScanSize {
		srverr.ResponseError(w, r, http.StatusBadRequest, errLimitRequestBodySizeMsg, h.logger)
		return
	}

	sr := http.MaxBytesReader(w, file, user.MaxScanSize)
	defer sr.Close()

	id, err := h.usecase.Upload(ctx, userID,
		model.Scan{
			DocumentID:  uint64(*sc.DocumentID),
			Type:        model.ScanType(sc.Type),
			Description: *sc.Description,
		},
		model.File{
			Reader:      sr,
			Size:        header.Size,
			ContentType: header.Header.Get("Content-Type"),
		})
	if err != nil {
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

	w.Header().Set("Location",
		api.BaseURL+"/users/"+strconv.FormatUint(userID, 10)+"/scans/"+strconv.FormatUint(id, 10))
}

func handleScanMultipartRequest(ctx context.Context, r *http.Request) (api.UploadScanMultipartRequestBody, error) {
	var scan api.UploadScanMultipartRequestBody

	scan.Type = api.ScanType(r.PostFormValue("type"))
	desc := r.PostFormValue("description")
	scan.Description = &desc
	err := scan.Validate(ctx, validator.Instance())
	if err != nil {
		return api.UploadScanMultipartRequestBody{}, err
	}

	var docID uint64
	if r.PostFormValue("document_id") != "" {
		docID, err = strconv.ParseUint(r.PostFormValue("document_id"), 10, 64)
		if err != nil {
			return api.UploadScanMultipartRequestBody{}, err
		}
		scan.DocumentID = &docID
	}

	return scan, nil
}

// @Router /users/{user_id}/scans/{scan_id} [delete]
func (h scanHandlers) DeleteScan(w http.ResponseWriter, r *http.Request, userID uint64, scanID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Scan
// @Router  /users/{user_id}/scans/{scan_id} [get]
func (h scanHandlers) GetScan(w http.ResponseWriter, r *http.Request, userID uint64, scanID uint64) {
	ctx := r.Context()

	scan, err := h.usecase.Get(ctx, userID, scanID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err, h.logger)
		return
	}

	if err = response.JSON(w, http.StatusOK, toAPIScan(scan)); err != nil {
		srverr.LogError(r, err, false, h.logger)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg,
			h.logger)
	}
}
