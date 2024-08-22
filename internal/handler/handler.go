package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Collap5e-dev/FlickPick/internal/config"
	_ "github.com/lib/pq"
)

type repo interface {
	GetMovieList(ctx context.Context) ([]MovieStruct, error)
}

func NewHandler(config *config.Config, repo repo) *Handler {
	return &Handler{
		config: config,
		repo:   repo,
	}
}

type Handler struct {
	config *config.Config
	repo   repo
}

type MovieStruct struct {
	Movie_id     int     `json:"movie_id"`
	Kinopoisk_id int     `json:"kinopoisk_id"`
	Year         int     `json:"year"`
	Name         string  `json:"name"`
	Genre        string  `json:"genre"`
	Rating_kp    float64 `json:"rating_kp"`
	Rating_imdb  float64 `json:"rating_imdb"`
	Rating_avg   float64 `json:"rating_avg"`
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	movieList, err := h.repo.GetMovieList(r.Context())
	if err != nil {
		panic(err)
	}
	body, err := json.Marshal(movieList)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, body)
}
