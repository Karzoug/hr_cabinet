package handlers

import (
	"net/http"
)

// @Produce application/json
// @Success 200 {array} api.Scan
// @Router  /users/{user_id}/scans [get]
func (h *handler) ListScans(w http.ResponseWriter, r *http.Request, userID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  multipart/form-data
// @Param   body body api.UploadScanMultipartRequestBody true ""
// @Router  /users/{user_id}/scans [post]
func (h *handler) UploadScan(w http.ResponseWriter, r *http.Request, userID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /users/{user_id}/scans/{scan_id} [delete]
func (h *handler) DeleteScan(w http.ResponseWriter, r *http.Request, userID int, scanID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Scan
// @Router  /users/{user_id}/scans/{scan_id} [get]
func (h *handler) GetScan(w http.ResponseWriter, r *http.Request, userID int, scanID int) {
	w.WriteHeader(http.StatusNotImplemented)
}
