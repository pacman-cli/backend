package main

import (
	"log"
	"net/http"

	"github.com/puspo/basicnewproject/internal/config"
	"github.com/puspo/basicnewproject/internal/db"
	"github.com/puspo/basicnewproject/internal/repositories"
	"github.com/puspo/basicnewproject/internal/services"
	posthandlers "github.com/puspo/basicnewproject/internal/transport/http/handlers"
	"github.com/puspo/basicnewproject/internal/transport/http/middleware"
)

func main() {
	// Load configuration from environment
	cfg := config.Load()

	// Initialize DB connection
	sqlDB, err := db.OpenMySQL(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer sqlDB.Close()

	// Wire dependencies across layers
	repo := repositories.NewPostRepository(sqlDB)
	service := services.NewPostService(repo)
	handlers := posthandlers.NewPostHandlers(service)

	// Set up routes
	mux := http.NewServeMux()
	handlers.Register(mux)

	// Apply middleware
	var handler http.Handler = mux
	handler = middleware.RecoveryMiddleware(handler)
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.CORS(handler)

	addr := ":" + cfg.AppPort
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
