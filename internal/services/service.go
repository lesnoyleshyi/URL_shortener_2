package services

import (
	"URL_shortener_2/internal/domain"
	"URL_shortener_2/internal/repository"
	"fmt"
	"time"
)

type Service interface {
	Save(longUrl string) (shortUrl string, err error)
	Get(shortUrl string) (*domain.Url, error)
}

type service struct {
	repo repository.Repository
}

func New(repo *repository.Repository) Service {
	return service{repo: *repo}
}

func (s service) Save(longUrl string) (string, error) {
	Url := domain.Url{
		LongUrl:   longUrl,
		CreatedAt: time.Now(),
	}
	Url.ShortUrl = genShort()
	if err := s.repo.Save(&Url); err != nil {
		return "", fmt.Errorf("unable to save url in sorage: %w", err)
	}
	return Url.ShortUrl, nil
}

func (s service) Get(url string) (*domain.Url, error) {
	return s.repo.Get(url)
}
