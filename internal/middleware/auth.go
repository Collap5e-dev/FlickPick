package middleware

import (
	"context"
	"net/http"

	"github.com/Collap5e-dev/FlickPick/internal/config"
)

func Auth(cfg *config.Config, method func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Auth logic
		// if auth_login is true
		ctx := context.WithValue(r.Context(), "username", "from jwt")
		r = r.WithContext(ctx)
		// then method(r, w)
		// else
		// не вызываем method(), а просто отдаешь 403 auth error
	}
}
