package main

import (
	// "fmt"
	"fmt"
	"net/http"
	"time"

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
	handler1 := handler.NewHandler(*cfg)
	router := mux.NewRouter()
	router.HandleFunc("/", handler1.Home)
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Running...")
	for {
		time.Sleep(time.Second)
	}

}
