package middleware

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Collap5e-dev/FlickPick/internal/config"
	"github.com/Collap5e-dev/FlickPick/internal/handler"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func Auth(cfg *config.Config, method func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "username", "from jwt")
		r = r.WithContext(ctx)
		h := &handler.Handler{}
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			h.HandlerError(w, 401, fmt.Errorf("токен не обнаружен"), "Токен не обнаружен")
			return
		}
		error := isValidToken(tokenString)
		if error != nil {
			h.HandlerError(w, 401, error, "Ошибка при проверке токена")
			return
		}

		token, err := VerifyToken(cfg, tokenString)
		if err != nil {
			h.HandlerError(w, 401, err, "Ошибка валидации токена")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			h.HandlerError(w, 401, fmt.Errorf("токен недействителен"), "Токен недействителен")
			return
		}
		fmt.Printf("Пользователь: %v\n", claims["username"])
		method(w, r)

	}
}

func VerifyToken(cfg *config.Config, tokenString string) (*jwt.Token, error) {
	secretKey := cfg.SecretKey
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		fmt.Println("Ошибка при парсинге токена:", err)
		return nil, err
	}

	return token, nil
}

func isValidToken(tokenString string) error {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return errors.New("token contains an invalid number of segments")
	}

	for i, part := range parts {
		switch i {
		case 0: // Header
			if !isValidBase64(part) {
				return errors.New("header segment is not valid base64")
			}
		case 1: // Payload
			if !isValidBase64(part) {
				return errors.New("payload segment is not valid base64")
			}
		case 2: // Signature
			if !isValidBase64(part) {
				return errors.New("signature segment is not valid base64")
			}
		default:
			return errors.New("unexpected part in token")
		}
	}

	return nil
}

func isValidBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
