package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Env menampung semua konfigurasi dari environment / .env
type Env struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret  string
	JWTExpireH int

	AppPort string
	AppEnv  string
}

func Load() *Env {
	_ = godotenv.Load() // .env opsional

	e := &Env{
		DBHost:     get("DB_HOST", "localhost"),
		DBPort:     get("DB_PORT", "5432"),
		DBUser:     get("DB_USER", "postgres"),
		DBPassword: get("DB_PASSWORD", ""),
		DBName:     get("DB_NAME", "siakad_mini"),
		DBSSLMode:  get("DB_SSLMODE", "disable"),
		JWTSecret:  get("JWT_SECRET", "rahasia-siakad-mini"),
		AppPort:    get("APP_PORT", "3000"),
		AppEnv:     get("APP_ENV", "development"),
	}
	e.JWTExpireH = 24
	if v := os.Getenv("JWT_EXPIRE_H"); v != "" {
		var h int
		_, _ = fmt.Sscanf(v, "%d", &h)
		if h > 0 {
			e.JWTExpireH = h
		}
	}
	return e
}

func (e *Env) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		e.DBHost, e.DBPort, e.DBUser, e.DBPassword, e.DBName, e.DBSSLMode)
}

func get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
