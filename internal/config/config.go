package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config menyimpan semua konfigurasi aplikasi yang diambil dari environment variables.
type Config struct {
	AppPort      string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	JWTSecret    string
	QRHMACSecret string
	RedisAddr    string
}

// Load membaca file .env (jika ada) lalu mengembalikan struct Config.
// Di production, environment variables biasanya di-inject langsung oleh orchestrator
// (Docker/K8s) sehingga file .env tidak wajib ada.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("info: file .env tidak ditemukan, menggunakan environment variables sistem")
	}

	return &Config{
		AppPort:      getEnv("APP_PORT", "8080"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "postgres"),
		DBName:       getEnv("DB_NAME", "coffeshop_db"),
		DBSSLMode:    getEnv("DB_SSLMODE", "disable"),
		JWTSecret:    getEnv("JWT_SECRET", "change-me-in-production"),
		QRHMACSecret: getEnv("QR_HMAC_SECRET", "change-me-too"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
