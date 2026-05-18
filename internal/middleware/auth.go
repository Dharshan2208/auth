package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(secret []byte, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(
			authHeader,
			"Bearer ",
		)

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return secret, nil
			},
		)

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		ctx := r.Context()
		if username, ok := claims["username"]; ok {
			ctx = context.WithValue(ctx, "username", username)
		}
		if email, ok := claims["email"]; ok {
			ctx = context.WithValue(ctx, "email", email)
		}
		if role, ok := claims["role"]; ok {
			ctx = context.WithValue(ctx, "role", role)
		}
		if rawID, ok := claims["user_id"]; ok {
			// jwt.MapClaims unmarshals numbers as float64.
			if f, ok := rawID.(float64); ok {
				ctx = context.WithValue(ctx, "user_id", int(f))
			} else {
				ctx = context.WithValue(ctx, "user_id", rawID)
			}
		}

		next(w, r.WithContext(ctx))
	}
}
