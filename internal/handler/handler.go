package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Collap5e-dev/FlickPick/internal/config"
	"github.com/Collap5e-dev/FlickPick/internal/model"
	_ "github.com/lib/pq"
)

type repo interface {
	GetMovieList(ctx context.Context) ([]model.Movie, error)
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
