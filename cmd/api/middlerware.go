package api

import (
	"context"
	"errors"
	"main/internal/evn"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func HandleJWTAuth(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeaderString := r.Header.Get("Authorization")

		if authHeaderString == "" {
			err := errors.New("authorization header required")
			unauthorized(w, r, err)
			return
		}

		tokenString := authHeaderString[7:]

		var secret_words string = "A request for a long text message: Search results showIf this is your intent, please clarify the context and what you want the text to be about."

		secret := evn.GetString(secret_words, "JWT_SERECT")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		})

		if err != nil {
			err := errors.New("invalid token")
			unauthorized(w, r, err)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		exp := claims["exp"]

		date_now := time.Now()

		t,ok := exp.(time.Time)

		if !ok {
			err := errors.New("token exp date cannot be confirmed")
			unauthorized(w, r, err)
			return
		}

		if t.Before(date_now) {
			err := errors.New("token expired")
			unauthorized(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		h.ServeHTTP(w, r.WithContext(ctx))
	})

}
