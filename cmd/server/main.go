package main

import (
	// "fmt"

	"fmt"
	"log"
	"net/http"
	"strconv"

	// "github.com/Collap5e-dev/FlickPick/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Collap5e-dev/FlickPick/internal/config"
	"github.com/Collap5e-dev/FlickPick/internal/handler"
	"github.com/Collap5e-dev/FlickPick/internal/repo"
	"github.com/gorilla/mux"
	// "time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	fmt.Printf("read config %v\n", cfg)
	psqlInfo := cfg.Db.Dsn()
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}
	defer db.Close()

	movieRepo := repo.NewMovieList(db)
	handler1 := handler.NewHandler(cfg, movieRepo)
	router := mux.NewRouter()
	router.HandleFunc("/", handler1.Home)
	http.Handle("/", router)
	http.ListenAndServe(":"+strconv.FormatInt(cfg.Port, 10), nil)
}
