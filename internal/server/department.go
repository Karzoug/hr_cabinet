package server

import "net/http"

// @Produce application/json
// @Success 200 {array} api.Department
// @Router  /departments [get]
func (s *server) ListDepartments(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
