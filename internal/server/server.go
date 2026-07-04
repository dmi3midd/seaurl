package server

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"seaurl/internal/config"
	"seaurl/internal/postgres"
	"seaurl/internal/repository"
	"seaurl/internal/service"
)

type Server struct {
	cfg     config.Config
	db      postgres.PostgresService
	service service.URLService
}

func NewServer(cfg *config.Config, db postgres.PostgresService) *http.Server {
	urlRepository := repository.NewURLStorage(db.GetDB())
	newServer := &Server{
		cfg:     *cfg,
		db:      db,
		service: service.NewURLService(urlRepository),
	}

	router := newServer.RegisterRoutes()
	// Declare Server config
	server := &http.Server{
		Addr:         cfg.Http.Address,
		Handler:      router,
		IdleTimeout:  cfg.Http.IdleTimeout,
		ReadTimeout:  cfg.Http.ReadTimeout,
		WriteTimeout: cfg.Http.WriteTimeout,
	}

	return server
}
