package main

import (
	"fmt"
	"net/http"

	"github.com/Collap5e-dev/FlickPick/internal/config"
	"github.com/Collap5e-dev/FlickPick/internal/handler"
	"github.com/gorilla/mux"

	// "github.com/Collap5e-dev/FlickPick/internal/handler"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()

	router.HandleFunc("/", handler.Home)
	http.Handle("/", router)
	http.ListenAndServe(":8080", router)
	fmt.Println(cfg)
	fmt.Println("Running...")
	for {
		time.Sleep(time.Second)
	}
}
