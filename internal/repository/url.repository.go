package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"seaurl/internal/models"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUrlNotFound = errors.New("url not found")
)

// URLRepository is an interface to interact with urls in db.
type URLRepository interface {
	// GetByAlias retrieves a Url entity by its alias.
	// It returns ErrUrlNotFound if no url are found.
	GetByAlias(ctx context.Context, alias string) *models.Url
	// Save creates a Url entity and returns it.
	Save(ctx context.Context, url string, alias string) *models.Url
}

type URLStorage struct {
	db *sqlx.DB
}

func NewURLStorage(db *sqlx.DB) URLStorage {
	return URLStorage{
		db: db,
	}
}

func (store *URLStorage) GetByAlias(ctx context.Context, alias string) (*models.Url, error) {
	op := "UrlRepository.GetByAlias"
	query := `SELECT id, url, alias FROM urls WHERE alias = ?`
	var url models.Url
	err := store.db.GetContext(ctx, &url, query, alias)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUrlNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &url, nil
}

func (store *URLStorage) Save(ctx context.Context, url *models.Url) (*models.Url, error) {
	op := "UrlRepository.Save"
	query := `INSERT INTO urls SET (id, url, alias) VALUES (:id, :url, :alias)`
	if _, err := store.db.NamedExecContext(ctx, query, url); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}
