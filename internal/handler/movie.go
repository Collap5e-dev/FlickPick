package handler

import (
	"net/http"
)

func (h *Handler) MovieAdd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := ctx.Value("username").(string)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(username))
}
