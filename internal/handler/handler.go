package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Collap5e-dev/FlickPick/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Handler struct {
	config *config.Config
	db     *sqlx.DB
}

type movieStruct struct {
	Movie_id     int     `json:"movie_id"`
	Kinopoisk_id int     `json:"kinopoisk_id"`
	Year         int     `json:"year"`
	Name         string  `json:"name"`
	Genre        string  `json:"genre"`
	Rating_kp    float64 `json:"rating_kp"`
	Rating_imdb  float64 `json:"rating_imdb"`
	Rating_avg   float64 `json:"rating_avg"`
}

func (h Handler) Home(w http.ResponseWriter, r *http.Request) {
	// data, err := json.MarshalIndent(h.config, "", "\t")
	// if err != nil {
	// 	panic(err)
	// }
	// formattedMessage := fmt.Sprintf("%s", data)
	// fmt.Fprint(w, formattedMessage)
	movieTable, err := h.db.Query(`
		SELECT
			movie_id,
			name,
			rating_kp,
			rating_imdb,
			kinopoisk_id,
			rating_avg,
			genre,
			year
		FROM
			content
		ORDER BY
			name
	`)
	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v\n", err)
	}
	defer movieTable.Close()
	movieList := make([][]byte, 0)
	for movieTable.Next() {
		var movie_id, kinopoisk_id, year int
		var name, genre string
		var rating_kp, rating_imdb, rating_avg float64
		err := movieTable.Scan(&movie_id, &name, &rating_kp, &rating_imdb, &kinopoisk_id, &rating_avg, &genre, &year)
		if err != nil {
			log.Fatalf("Ошибка при сканировании строки: %v\n", err)
			continue
		}
		mStruct := movieStruct{
			Movie_id:     movie_id,
			Name:         name,
			Rating_kp:    rating_kp,
			Rating_imdb:  rating_imdb,
			Kinopoisk_id: kinopoisk_id,
			Rating_avg:   rating_avg,
			Genre:        genre,
			Year:         year}
		data, err := json.MarshalIndent(mStruct, "", "\t")
		if err != nil {
			panic(err)
		}
		movieList = append(movieList, data)
	}
	formattedMessage := fmt.Sprintf("%s", movieList)
	fmt.Fprint(w, formattedMessage)
}

func NewHandler(config *config.Config, db *sqlx.DB) *Handler {
	return &Handler{
		config: config,
		db:     db,
	}
}
