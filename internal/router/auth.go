package router

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatihsezgin/candlecloud-backend/internal/app"
	"github.com/urfave/negroni"
)

func Auth() negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var tokenstr string
		bearerToken := r.Header.Get("Authorization")
		strArr := strings.Split(bearerToken, " ")
		if len(strArr) == 2 {
			tokenstr = strArr[1]
		}

		token, err := app.TokenValid(tokenstr)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)

		// Get User UUID from claims
		ctxUserUUID, ok := claims["user_uuid"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Admin or Member
		ctxAuthorized, ok := claims["authorized"].(bool)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctxWithUUID := context.WithValue(ctx, "uuid", ctxUserUUID)
		ctxWithAuthorized := context.WithValue(ctxWithUUID, "authorized", ctxAuthorized)
		next(w, r.WithContext(ctxWithAuthorized))
	})
}
