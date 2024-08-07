package handler

import (
	"fmt"
	"net/http"

	"github.com/Collap5e-dev/FlickPick/internal/config"
)

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome Home!")
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
