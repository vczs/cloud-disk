package middleware

import (
	"cloud-disk/disk/helper"
	"net/http"
)

type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) AuthMiddlewareHandle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("token is empty!"))
			return
		}
		user, err := helper.AuthToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		r.Header.Set("Id", string(rune(user.Id)))
		r.Header.Set("Uid", user.Uid)
		r.Header.Set("Name", user.Name)
		next(w, r)
	}
}
