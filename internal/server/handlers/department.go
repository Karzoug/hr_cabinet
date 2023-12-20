package handlers

import (
	"net/http"
)

// @Produce application/json
// @Success 200 {array} api.Department
// @Router  /departments [get]
func (h *handler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
