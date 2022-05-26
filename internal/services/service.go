package services

import (
	"URL_shortener_2/internal/domain"
	"URL_shortener_2/internal/repository"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

//go:generate mockgen -source=service.go -destination=moks/mock.go

type Service interface {
	//Save saves given longUrl, generates and returns shortUrl or error
	Save(longUrl string) (shortUrl string, err error)
	//Get receives url and returns it's "opposite" url:
	//if shortUrl is given, it returns longUrl and vice versa
	Get(shortOrLongUrl string) (longOrShortUrl string, err error)
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

func (s service) Get(url string) (string, error) {
	urlStruct, err := s.repo.Get(url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repository.ErrNoSuchUrl
		}
		return "", err
	}
	//As we don't divide Get method to 'getByShort' and 'getByLong',
	//we have to decide what url should be returned: long or short.
	if urlStruct.LongUrl == url {
		return urlStruct.ShortUrl, nil
	} else {
		return urlStruct.LongUrl, nil
	}
}
