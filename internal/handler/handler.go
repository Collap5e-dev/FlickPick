package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"net/http"
	"net/mail"
	"time"

	_ "github.com/lib/pq"

	"github.com/Collap5e-dev/FlickPick/internal/config"
	"github.com/Collap5e-dev/FlickPick/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type repo interface {
	GetMovieList(ctx context.Context) ([]model.Movie, error)
	CreateUser(ctx context.Context, user model.User) error
	GiveUserPass(ctx context.Context, loginData string) (string, error)
	CreateNewMovie(ctx context.Context, NewMovie model.Movie) error
}

func NewHandler(config *config.Config, repo repo) *Handler {
	return &Handler{
		config: config,
		repo:   repo,
	}
}

type TokenResponse struct {
	Token string `json:"token"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Handler handle requests
type Handler struct {
	config *config.Config
	repo   repo
}

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
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
	w.WriteHeader(http.StatusOK)
	w.Write(body)

}

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	var userData model.User
	body, err := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &userData); err != nil {
		h.HandlerError(w, 500, err, "ошибка обработки данных")
		return
	}
	defer r.Body.Close()
	err, valid := validData(userData)
	if err != nil {
		h.HandlerError(w, 400, err, "ошибка валидации данных")
		return
	}
	if valid {
		fmt.Println("Данные валидны")
	}
	username := userData.Username
	password, err := hashPassword(userData.Password)
	if err != nil {
		h.HandlerError(w, 500, err, "ошибка хэширования пароля")
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
		h.HandlerError(w, 500, errCreateUser, "ошибка создания пользователя")
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var userData LoginUser
	body, err := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &userData); err != nil {
		h.HandlerError(w, 500, err, "ошибка обработки данных")
		return
	}
	if err != nil {
		h.HandlerError(w, 500, err, "ошибка чтения данных")
		return
	}
	defer r.Body.Close()

	pass, err := h.repo.GiveUserPass(r.Context(), userData.Username)
	if err != nil {
		h.HandlerError(w, 500, err, "ошибка проверки пользователя")
		return
	}
	if !checkHashPassword(userData.Password, pass) {
		h.HandlerError(w, 400, err, "неверно введен пароль")
		return
	}
	token, err := h.createToken(userData.Username)
	if err != nil {
		h.HandlerError(w, 500, err, "ошибка создания токена")
		return
	}
	sendToken, err := json.Marshal(token)
	if err != nil {
		h.HandlerError(w, 500, err, "ошибка отправки токена")
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(sendToken)

}

func (h *Handler) HandlerError(w http.ResponseWriter, statusCode int, err error, text string) {
	fmt.Println(text, err)
	m := Error{
		Status:  statusCode,
		Message: text,
	}
	message, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(message)
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

func (h *Handler) createToken(username string) (TokenResponse, error) {
	var newToken TokenResponse
	secretKey := h.config.SecretKey
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 168).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return newToken, err
	}
	newToken = TokenResponse{
		Token: tokenString,
	}
	return newToken, nil
}
