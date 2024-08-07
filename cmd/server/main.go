package main

import (
	"fmt"

	"github.com/Collap5e-dev/FlickPick/internal/config"
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
	//router.Get("/", handler.Home)
	//http.ListenAndServe(":8080", router)
	fmt.Println(cfg)
	fmt.Println("Running...")
	for {
		time.Sleep(time.Second)
	}
}
