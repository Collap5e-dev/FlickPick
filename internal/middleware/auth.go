package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func Auth(secretKey []byte, method func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			HandlerError(w, 401, fmt.Errorf("токен не обнаружен"), "Токен не обнаружен")
			return
		}

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		}, jwt.WithExpirationRequired())
		if err != nil {
			HandlerError(w, 401, err, "Ошибка валидации токена")
			return
		}
		if !token.Valid {
			HandlerError(w, 401, fmt.Errorf("токен недействителен"), "Токен недействителен")
			return
		}
		cl := *claims
		ctx := r.Context()
		ctx = context.WithValue(ctx, "username", cl["username"])
		r = r.WithContext(ctx)
		method(w, r)

	}
}

func HandlerError(w http.ResponseWriter, statusCode int, err error, text string) {
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
