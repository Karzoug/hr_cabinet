package handlers

import (
	"net/http"
)

// @Success 200
// @Router  /health [get]
func (h *handler) Health(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()

	// TODO: implement health check
	// ping DB may be useful

	w.WriteHeader(http.StatusOK)
}
