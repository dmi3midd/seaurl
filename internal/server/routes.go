package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"seaurl/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	router := chi.NewRouter()

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
			slog.Error(
				"failed to decode request body",
				slog.String("error", err.Error()),
				slog.String("url", r.URL.String()),
			)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		url, err := s.service.Save(r.Context(), reqBody.Url)
		if err != nil {
			slog.Error(
				"failed to create url",
				slog.String("error", err.Error()),
				slog.String("url", reqBody.Url),
			)
			http.Error(w, "Failed to save URL", http.StatusInternalServerError)
			return
		}

		slog.Info(
			"url created successfully",
			slog.String("url", url.Url),
			slog.String("alias", url.Alias),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(url); err != nil {
			slog.Error(
				"failed to encode response body",
				slog.String("error", err.Error()),
				slog.String("url", r.URL.String()),
			)
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
				slog.Info(
					"url not found",
					slog.String("error", err.Error()),
					slog.String("alias", alias),
				)
				http.Error(w, "URL not found", http.StatusNotFound)
				return
			}
			slog.Error(
				"failed to get url",
				slog.String("error", err.Error()),
				slog.String("alias", alias),
			)
			http.Error(w, "Failed to get URL", http.StatusInternalServerError)
			return
		}

		slog.Info(
			"redirected successfully",
			slog.String("alias", url.Alias),
			slog.String("to", url.Url),
		)
		http.Redirect(w, r, url.Url, http.StatusMovedPermanently)
	}
}
