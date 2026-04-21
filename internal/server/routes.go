package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"seaurl/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Post("/", s.CreateAliasHandler())
	router.Get("/{alias}", s.RedirectHandler())

	return router
}

type ReqBody struct {
	Url string `json:"url"`
}

func (s *Server) CreateAliasHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody ReqBody
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		url, err := s.service.Save(r.Context(), reqBody.Url)
		if err != nil {
			http.Error(w, "Failed to save URL", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(url); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	}
}

func (s *Server) RedirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}

		url, err := s.service.GetByAlias(r.Context(), alias)
		if err != nil {
			if errors.Is(err, service.ErrUrlNotFound) {
				http.Error(w, "URL not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to get URL", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url.Url, http.StatusMovedPermanently)
	}
}
