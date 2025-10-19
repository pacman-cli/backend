// Package config handles environment-based configuration.
package config

import (
	"net/url"
	"os"
	"strconv"
	"strings"
)

// AppConfig holds all configuration values.
type AppConfig struct {
	AppPort       string
	DBUser        string
	DBPass        string
	DBHost        string
	DBPort        string
	DBName        string
	RunMigrations bool
}

// Load reads environment variables and applies sensible defaults.
// Supports either DATABASE_URL or individual DB_* variables.
func Load() AppConfig {
	appPort := getEnv("PORT", "8080")
	runMigrations := getEnvBool("RUN_MIGRATIONS", false)

	// Prefer DATABASE_URL (Railway, Render, etc.)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		user, pass, host, port, name := parseMySQLURL(dbURL)
		return AppConfig{
			AppPort:       appPort,
			DBUser:        user,
			DBPass:        pass,
			DBHost:        host,
			DBPort:        port,
			DBName:        name,
			RunMigrations: runMigrations,
		}
	}

	// Fallback to local dev variables
	return AppConfig{
		AppPort:       appPort,
		DBUser:        getEnv("DB_USER", "root"),
		DBPass:        getEnv("DB_PASS", ""),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "3306"),
		DBName:        getEnv("DB_NAME", "blogdbGoLang"),
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

// parseMySQLURL parses mysql://user:pass@host:port/dbname into its components.
func parseMySQLURL(mysqlURL string) (user, pass, host, port, name string) {
	u, err := url.Parse(mysqlURL)
	if err != nil {
		panic("Invalid DATABASE_URL format: " + err.Error())
	}

	user = u.User.Username()
	pass, _ = u.User.Password()
	host = u.Hostname()
	port = u.Port()
	name = strings.TrimPrefix(u.Path, "/")

	if port == "" {
		port = "3306"
	}
	return user, pass, host, port, name
}
