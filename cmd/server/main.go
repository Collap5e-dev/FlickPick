package main

import (
	// "fmt"

	"fmt"
	"net/http"
	"strconv"

	// "github.com/Collap5e-dev/FlickPick/internal/config"

	"github.com/Collap5e-dev/FlickPick/internal/config"
	"github.com/Collap5e-dev/FlickPick/internal/handler"
	"github.com/gorilla/mux"
	// "time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	fmt.Printf("read config %v\n", cfg)
	handler1 := handler.NewHandler(*cfg)
	router := mux.NewRouter()
	router.HandleFunc("/", handler1.Home)
	http.Handle("/", router)
	http.ListenAndServe(":"+strconv.FormatInt(cfg.Port, 10), nil)
}
