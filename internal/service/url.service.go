package service

import (
	"context"
	"errors"
	"fmt"
	"seaurl/internal/models"
	"seaurl/internal/repository"

	"github.com/rs/xid"
)

var (
	ErrUrlNotFound = errors.New("url not found")
)

type URLService interface {
	// GetByAlias retunrs a Url by its alias.
	// It returns ErrUrlNotFound if no url are found.
	GetByAlias(ctx context.Context, alias string) (*models.Url, error)
	// Save saves Url and returns it.
	Save(ctx context.Context, url string) (*models.Url, error)
}

type service struct {
	store repository.URLRepository
}

func NewURLService(store repository.URLRepository) *service {
	return &service{
		store: store,
	}
}

func (s *service) GetByAlias(ctx context.Context, alias string) (*models.Url, error) {
	op := "URLService.GetByAlias"
	url, err := s.store.GetByAlias(ctx, alias)
	if err != nil {
		if errors.Is(err, repository.ErrUrlNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrUrlNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}

// TODO: Add collision check. Currently relying on high entropy of crypto/rand.
// TODO: Add validation for url.
func (s *service) Save(ctx context.Context, urlStr string) (*models.Url, error) {
	op := "URLService.Save"

	id := xid.New().String()[:10]
	alias := xid.New().String()[:8]
	newUrl := models.Url{
		Id:    id,
		Url:   urlStr,
		Alias: alias,
	}

	resUrl, err := s.store.Create(ctx, &newUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resUrl, nil
}
