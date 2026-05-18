package auth

import (
	"time"

	"github.com/Dharshan2208/auth/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(user models.User, secret []byte, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(secret)
}

func GenerateRefreshToken(user models.User, secret []byte, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"type":     "refresh",
		"exp":      time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(secret)
}
