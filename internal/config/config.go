package config

import (
	"os"
	"strconv"
)

// AppConfig holds all configuration values loaded from environment variables.
// The intent is to keep configuration centralized and explicit.
type AppConfig struct {
	AppPort       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPass        string
	DBName        string
	RunMigrations bool
}

// Load reads environment variables and applies sensible defaults for local dev.
// No external libraries are used to keep the workflow transparent.
func Load() AppConfig {
	appPort := getEnv("APP_PORT", "8080")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "MdAshikur123+")
	dbName := getEnv("DB_NAME", "blogdbGoLang")
	runMigrations := getEnvBool("RUN_MIGRATIONS", false)

	return AppConfig{
		AppPort:       appPort,
		DBHost:        dbHost,
		DBPort:        dbPort,
		DBUser:        dbUser,
		DBPass:        dbPass,
		DBName:        dbName,
		RunMigrations: runMigrations,
	}
}

// getEnv returns the environment value if set, otherwise the default.
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// getEnvBool parses a boolean environment variable with a default fallback.
func getEnvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return def
}
