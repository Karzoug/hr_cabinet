package handler

import (
	"net/http"
)

// @Produce application/json
// @Success 200 {object} api.ListDepartmentsJSONRequestBody
// @Router  /departments [get]
func (h *handler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
