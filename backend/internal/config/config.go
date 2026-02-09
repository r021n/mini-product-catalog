package config

import (
	"os"
	"strings"
)

type Config struct {
	Port           string
	DBURL          string
	JWTSecret      string
	AllowedOrigins []string
}

func Load() Config {
	port := getenv("PORT", "8080")
	dbURL := getenv("DB_URL", "postgres://postgres:postgres@localhost:5432/catalog?sslmode=disable")
	jwt := getenv("JWT_SECRET", "dev-secret-change-me")

	originsRaw := getenv("CORS_ORIGINS", "http://localhost:5173")
	origins := splitAndTrim(originsRaw)

	return Config{
		Port:           port,
		DBURL:          dbURL,
		JWTSecret:      jwt,
		AllowedOrigins: origins,
	}
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	return v
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}

	return out
}
