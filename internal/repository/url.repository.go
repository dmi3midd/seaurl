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
	GetByAlias(ctx context.Context, alias string) (*models.Url, error)
	// Create creates a Url entity and returns it.
	Create(ctx context.Context, url *models.Url) (*models.Url, error)
}

type repository struct {
	db *sqlx.DB
}

func NewURLStorage(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}

func (store *repository) GetByAlias(ctx context.Context, alias string) (*models.Url, error) {
	op := "URLRepository.GetByAlias"
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

func (store *repository) Create(ctx context.Context, url *models.Url) (*models.Url, error) {
	op := "URLRepository.Create"
	query := `INSERT INTO urls (id, url, alias) VALUES (:id, :url, :alias)`
	if _, err := store.db.NamedExecContext(ctx, query, url); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}
