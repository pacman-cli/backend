package main

import (
	"database/sql"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/go-sql-driver/mysql"

	"github.com/puspo/basicnewproject/internal/config"
)

// This is a very small migration runner for local development.
// It loads SQL files in the migrations directory and applies them in order.
func main() {
	cfg := config.Load()

	dsn := cfg.DBUser + ":" + cfg.DBPass + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci&multiStatements=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	migrationsDir := "migrations"
	var files []string
	if err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".sql" {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		log.Fatalf("list migrations: %v", err)
	}
	sort.Strings(files)

	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			log.Fatalf("read %s: %v", f, err)
		}
		log.Printf("applying %s", f)
		if _, err := db.Exec(string(b)); err != nil {
			log.Fatalf("apply %s: %v", f, err)
		}
	}
	log.Println("migrations applied")
}
