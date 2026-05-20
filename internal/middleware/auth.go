package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func Auth(secret []byte, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			writeUnauthorized(w, "missing authorization token")
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
			writeUnauthorized(w, "invalid or expired token")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if t, ok := claims["type"].(string); ok && t == "refresh" {
			writeUnauthorized(w, "wrong token type")
			return
		}

		ctx := r.Context()
		if role, ok := claims["role"]; ok {
			ctx = context.WithValue(ctx, "role", role)
		}
		if rawID, ok := claims["user_id"]; ok {
			if f, ok := rawID.(float64); ok {
				ctx = context.WithValue(ctx, "user_id", int(f))
			} else {
				ctx = context.WithValue(ctx, "user_id", rawID)
			}
		}

		next(w, r.WithContext(ctx))
	}
}
