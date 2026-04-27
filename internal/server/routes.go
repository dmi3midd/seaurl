package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	errs "seaurl/internal/errors"
	"seaurl/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	ErrEmptyAlias = errors.New("empty alias")
)

func (s *Server) RegisterRoutes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	router.Post("/", errs.ErrorHandler(s.CreateAliasHandler))
	router.Get("/{alias}", errs.ErrorHandler(s.RedirectHandler))

	return router
}

type ReqBody struct {
	Url string `json:"url"`
}

func (s *Server) CreateAliasHandler(w http.ResponseWriter, r *http.Request) error {
	var reqBody ReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		slog.Error(
			"failed to decode request body",
			slog.String("error", err.Error()),
			slog.String("url", r.URL.String()),
		)
		// return errs.InternalServerError(err)
	}

	url, err := s.service.Save(r.Context(), reqBody.Url)
	if err != nil {
		return errs.InternalServerError(err)
	}

	slog.Info(
		"alias created successfully",
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
		// return errs.InternalServerError(err)
	}

	return nil
}

func (s *Server) RedirectHandler(w http.ResponseWriter, r *http.Request) error {
	alias := chi.URLParam(r, "alias")
	if alias == "" {
		// http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return errs.NewBadRequestError(ErrEmptyAlias, "Empty alias")
	}

	url, err := s.service.GetByAlias(r.Context(), alias)
	if err != nil {
		if errors.Is(err, service.ErrUrlNotFound) {
			// slog.Info(
			// 	"url not found",
			// 	slog.String("error", err.Error()),
			// 	slog.String("alias", alias),
			// )
			return errs.NewNotFoundError(err, "page not found")
		}
		return errs.InternalServerError(err)
	}

	slog.Info(
		"redirected successfully",
		slog.String("alias", url.Alias),
		slog.String("to", url.Url),
	)
	http.Redirect(w, r, url.Url, http.StatusMovedPermanently)

	return nil
}
