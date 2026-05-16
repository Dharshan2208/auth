package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     mustEnv("DATABASE_URL"),
		JWTSecret:       mustEnv("JWT_SECRET"),
		AccessTokenTTL:  mustDuration("ACCESS_TOKEN_TTL", "10m"),
		RefreshTokenTTL: mustDuration("REFRESH_TOKEN_TTL", "168h"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env: %s", key)
	}
	return val
}

func mustDuration(key string, fallback string) time.Duration {
	val := getEnv(key, fallback)

	d, err := time.ParseDuration(val)
	if err != nil {
		log.Fatalf("invalid duration for %s", key)
	}

	return d
}
