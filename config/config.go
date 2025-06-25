package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DB_DSN     string
	JWTSecret  string
	DB_CA_CERT string
}

func LoadConfig() *Config {
	cfg := &Config{
		DB_DSN:     os.Getenv("DB_DSN"),
		JWTSecret:  getEnv("JWT_SECRET", "default_secret"),
		DB_CA_CERT: getEnv("DB_CA_CERT", "ca.pem"),
	}

	// Convert relative paths to absolute
	if !filepath.IsAbs(cfg.DB_CA_CERT) {
		wd, _ := os.Getwd()
		cfg.DB_CA_CERT = filepath.Join(wd, cfg.DB_CA_CERT)
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
