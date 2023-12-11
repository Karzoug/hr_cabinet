package server

import "net/http"

// @Produce application/json
// @Success 200 {array} api.Scan
// @Router  /users/{user_id}/scans [get]
func (s *server) ListScans(w http.ResponseWriter, r *http.Request, userID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Accept  multipart/form-data
// @Param   body body api.UploadScanMultipartRequestBody true ""
// @Router  /users/{user_id}/scans [post]
func (s *server) UploadScan(w http.ResponseWriter, r *http.Request, userID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Router /users/{user_id}/scans/{scan_id} [delete]
func (s *server) DeleteScan(w http.ResponseWriter, r *http.Request, userID int, scanID int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// @Produce application/json
// @Success 200 {object} api.Scan
// @Router  /users/{user_id}/scans/{scan_id} [get]
func (s *server) GetScan(w http.ResponseWriter, r *http.Request, userID int, scanID int) {
	w.WriteHeader(http.StatusNotImplemented)
}
