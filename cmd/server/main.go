package main

import (
	"fmt"
	"github.com/oke11o/abrakadabra/internal/config"

	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	//router := mux.CreateRouter
	//router.Get("/", handler.Home)
	//http.ListenAndServe(":8080", router)
	fmt.Println(cfg)
	fmt.Println("Running...")
	for {
		time.Sleep(time.Second)
	}
}
