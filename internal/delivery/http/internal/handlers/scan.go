package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/muonsoft/validation/validator"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/convert"
	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
	uservice "github.com/Employee-s-file-cabinet/backend/internal/service/user"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

// @Produce application/json
// @Success 200 {array} api.Scan
// @Router  /users/{user_id}/scans [get]
func (h *handler) ListScans(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	scans, err := h.userService.ListScans(ctx, userID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	if err = response.JSON(w, http.StatusOK, convert.ToAPIScans(scans)); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
	}
}

// @Accept  multipart/form-data
// @Param   body body api.UploadScanMultipartRequestBody true ""
// @Router  /users/{user_id}/scans [post]
func (h *handler) UploadScan(w http.ResponseWriter, r *http.Request, userID uint64) {
	ctx := r.Context()

	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	scan, err := handleScanMultipartRequest(ctx, r)
	if err != nil {
		srverr.ResponseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	file, header, err := r.FormFile("fileName")
	if err != nil {
		srverr.ResponseError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if header.Size > uservice.MaxScanSize {
		srverr.ResponseError(w, r, http.StatusBadRequest, errLimitRequestBodySizeMsg)
		return
	}

	sr := http.MaxBytesReader(w, file, uservice.MaxScanSize)
	defer sr.Close()

	id, err := h.userService.UploadScan(ctx, userID,
		model.Scan{
			DocumentID:  uint64(*scan.DocumentID),
			Type:        model.ScanType(scan.Type),
			Description: *scan.Description,
		},
		model.File{
			Reader:      sr,
			Size:        header.Size,
			ContentType: header.Header.Get("Content-Type"),
		})
	if err != nil {
		if errors.Is(err, new(http.MaxBytesError)) {
			srverr.ResponseError(w, r, http.StatusBadRequest, errLimitRequestBodySizeMsg)
			return
		}
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
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
func (h *handler) DeleteScan(w http.ResponseWriter, r *http.Request, userID uint64, scanID uint64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Scan
// @Router  /users/{user_id}/scans/{scan_id} [get]
func (h *handler) GetScan(w http.ResponseWriter, r *http.Request, userID uint64, scanID uint64) {
	ctx := r.Context()

	scan, err := h.userService.GetScan(ctx, userID, scanID)
	if err != nil {
		srverr.ResponseServiceError(w, r, err)
		return
	}

	if err = response.JSON(w, http.StatusOK, convert.ToAPIScan(scan)); err != nil {
		srverr.LogError(r, err, false)
		srverr.ResponseError(w, r,
			http.StatusInternalServerError,
			srverr.ErrInternalServerErrorMsg)
	}
}
