package auth

import (
	"time"

	"github.com/Dharshan2208/auth/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(user models.User, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(10 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(secret)
}

func GenerateRefreshToken(user models.User, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"type":     "refresh",
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(secret)
}
