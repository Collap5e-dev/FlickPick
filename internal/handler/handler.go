package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Collap5e-dev/FlickPick/internal/config"
)

type Handler struct {
	cfg config.Config
}

func Home(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		panic(err)
	}
	formattedMessage := fmt.Sprintf("%s", data)
	fmt.Fprint(w, formattedMessage)
}
