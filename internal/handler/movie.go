package handler

import (
	"net/http"
)

func (h *Handler) MovieAdd(w http.ResponseWriter, r *http.Request) {
	test := []byte("test")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(test)
}
