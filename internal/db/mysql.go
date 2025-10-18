package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/puspo/basicnewproject/internal/config"
)

// OpenMySQL opens a MySQL connection using database/sql and returns the *sql.DB pool.
// The connection string (DSN) is constructed from config. We also configure
// connection pool parameters appropriate for typical small services.
func OpenMySQL(cfg config.AppConfig) (*sql.DB, error) {
	// DSN format: user:pass@tcp(host:port)/dbname?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Tune the pool conservatively for development; adjust in production as needed.
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Verify the connection up front so we fail fast.
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
