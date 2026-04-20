package server

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"seaurl/internal/config"
	"seaurl/internal/database"
)

type Server struct {
	cfg config.Config
	db  database.DBService
}

func NewServer(cfg *config.Config, db *database.DBService) *http.Server {

	// Declare Server config
	server := &http.Server{
		Addr:         cfg.Http.Address,
		Handler:      http.NewServeMux(),
		IdleTimeout:  cfg.Http.IdleTimeout,
		ReadTimeout:  cfg.Http.ReadTimeout,
		WriteTimeout: cfg.Http.WriteTimeout,
	}

	return server
}
