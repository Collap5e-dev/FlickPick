package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/Collap5e-dev/FlickPick/internal/config"
	"github.com/Collap5e-dev/FlickPick/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type repo interface {
	GetMovieList(ctx context.Context) ([]model.Movie, error)
	CreateUser(ctx context.Context, user model.User) error
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

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	//_ = ctx
	var userData model.User
	body, err := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &userData); err != nil {
		fmt.Errorf("ошибка обработки данных: %w", err)
		return
	}
	defer r.Body.Close()
	username := userData.Username
	password := userData.Password
	email := userData.Email
	newUser := model.User{
		Username: username,
		Password: password,
		Email:    email,
	}
	errCreateUser := h.repo.CreateUser(r.Context(), newUser)
	if errCreateUser != nil {
		fmt.Errorf("ошибка создания пользователя: %w", err)
		return
	}
	w.WriteHeader(http.StatusOK)

	/*
		1 Получить из r *http.Request username, password
		2 Создать user := model.User{username: username, password: password}
		3 Вызвать метод repo.CreateUser(ctx, user)
		4 Написать в w http.ResponseWriter ответ типа {"Status":"OK"}


	*/
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: implement me
	//ctx := r.Context()
	//_ = ctx

}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkHashPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
