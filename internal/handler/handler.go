package handler

import (
	"context"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/Collap5e-dev/FlickPick/internal/config"
	"github.com/Collap5e-dev/FlickPick/internal/model"
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

// Handler handle requests
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
	w.Header().Add("Content-Type", "application/json")
	w.Write(body)
	w.WriteHeader(http.StatusOK)
}
