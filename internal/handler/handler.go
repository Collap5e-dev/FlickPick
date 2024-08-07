package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Collap5e-dev/FlickPick/internal/config"
)

type Handler struct {
	config config.Config
}

func (h Handler) Home(w http.ResponseWriter, r *http.Request) {
	data, err := json.MarshalIndent(h.config, "", "\t")
	if err != nil {
		panic(err)
	}
	formattedMessage := fmt.Sprintf("%s", data)
	fmt.Fprint(w, formattedMessage)
}

func NewHandler(config config.Config) *Handler {
	return &Handler{config: config}
}
