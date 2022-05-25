package repository

import (
	"URL_shortener_2/internal/domain"
	"errors"
)

type Repository interface {
	Save(url *domain.Url) error
	Get(Url string) (*domain.Url, error)
}

var ErrNoSuchUrl = errors.New("no such url is in storage")

func New(storageType string) Repository {
	if storageType == "cache" {
		return newCache()
	} else {
		return newDbRepo()
	}
}
