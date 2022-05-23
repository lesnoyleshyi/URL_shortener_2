package services

import "URL_shortener_2/internal/repository"

type Service interface {
	Save(shortUrl string)
	Get(longUrl string) string
}

type service struct {
	repo repository.Repository
}

func New(repo *repository.Repository) Service {
	return service{repo: *repo}
}

func (s service) Save(longUrl string) {
	s.repo.Save(longUrl)
}

func (s service) Get(shortUrl string) string {
	return s.repo.Get(shortUrl)
}
