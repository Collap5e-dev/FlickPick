package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) MovieAdd(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(map[string]string{"list": "empty"})
	if err != nil {
		h.handlerError(w, 500, err, "ошибка отправки токена")
		return
	}
	w.WriteHeader(http.StatusOK)
}
