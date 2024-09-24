package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"time"

	_ "github.com/lib/pq"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Collap5e-dev/FlickPick/internal/config"
	"github.com/Collap5e-dev/FlickPick/internal/model"
)

type repo interface {
	GetMovieList(ctx context.Context) ([]model.Movie, error)
	CreateUser(ctx context.Context, user model.User) error
	GiveUserPass(ctx context.Context, loginData string) (string, error)
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
		h.handlerError(w, 500, err, "ошибка обработки данных")
		return
	}
	defer r.Body.Close()
	err, valid := validData(userData)
	if err != nil {
		h.handlerError(w, 400, err, "ошибка валидации данных")
		return
	}
	if valid {
		fmt.Println("Данные валидны")
	}
	username := userData.Username
	password, err := hashPassword(userData.Password)
	if err != nil {
		h.handlerError(w, 500, err, "ошибка хэширования пароля")
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
		h.handlerError(w, 500, errCreateUser, "ошибка создания пользователя")
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
	var userData model.User
	body, err := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &userData); err != nil {
		h.handlerError(w, 500, err, "ошибка обработки данных")
		return
	}
	if err != nil {
		h.handlerError(w, 500, err, "ошибка чтения данных")
		return
	}
	defer r.Body.Close()

	pass, err := h.repo.GiveUserPass(r.Context(), userData.Username)
	if err != nil {
		h.handlerError(w, 500, err, "ошибка проверки пользователя")
		return
	}
	if !checkHashPassword(userData.Password, pass) {
		h.handlerError(w, 400, err, "неверно введен пароль")
		return
	}
	token, err := h.createToken(userData.Username)
	if err != nil {
		h.handlerError(w, 500, err, "ошибка создания токена")
		return
	}
	errorEncode := json.NewEncoder(w).Encode(map[string]string{"token": token})
	if errorEncode != nil {
		h.handlerError(w, 500, err, "ошибка отправки токена")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handlerError(w http.ResponseWriter, statusCode int, err error, text string) {
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

func (h *Handler) createToken(username string) (string, error) {
	secretKey := h.config.SecretKey
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 168).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//func verifyToken(tokenStr string) (*jwt.Token, error) {
//	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return secretKey, nil
//	})
//	if err != nil {
//		return nil, err
//	}
//	return token, nil
//}
