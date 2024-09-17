package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/mail"

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
	var userData model.User
	body, err := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &userData); err != nil {
		fmt.Errorf("ошибка обработки данных: %w", err)
		return
	}
	defer r.Body.Close()
	err, valid := validData(userData)
	if err != nil {
		fmt.Errorf("ошибка валидации данных: %w", err)
		return
	}
	if valid {
		fmt.Print("Данные валидны, регистрация прошла успешно")
	}
	username := userData.Username
	password, err := hashPassword(userData.Password)
	if err != nil {
		fmt.Errorf("ошибка хэширования пароля: %w", err)
		return
	}
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

func validData(user model.User) (error, bool) {
	_, errMail := mail.ParseAddress(user.Email)
	if len(user.Password) < 8 {
		err := errors.New("пароль недостаточно надежный")
		return err, false
	} else if len(user.Username) < 3 {
		err := errors.New("логин недостаточно надежный")
		return err, false
	} else if errMail != nil {
		return errMail, false
	} else {
		return nil, true
	}
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
