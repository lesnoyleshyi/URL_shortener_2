package repository

import (
	"URL_shortener_2/internal/domain"
	"errors"
	"sync"
)

//cacheRepo contains pointers to two maps:
//mapShortKey is a map[shortUrl]*domain.Url and
//mapLongKey is a map[longUrl]*domain.Url.
//Cache will contain a lot of data thus we need
//O(1) for both: insert in storage and retrieving from it.
//It's
type cacheRepo struct {
	sync.RWMutex
	mapShortKey *map[string]*domain.Url
	mapLongKey  *map[string]*domain.Url
}

func newCache() *cacheRepo {
	mapS := make(map[string]*domain.Url)
	mapL := make(map[string]*domain.Url)
	cacheStruct := cacheRepo{
		mapShortKey: &mapS,
		mapLongKey:  &mapL,
	}
	return &cacheStruct
}

func (c *cacheRepo) Save(url *domain.Url) error {
	c.Lock()
	(*c.mapLongKey)[url.LongUrl] = url
	(*c.mapShortKey)[url.ShortUrl] = url
	c.Unlock()
	return nil
}

func (c *cacheRepo) GetByShort(shortUrl string) (*domain.Url, error) {
	c.RLock()
	url, ok := (*c.mapShortKey)[shortUrl]
	c.RUnlock()
	if !ok {
		return nil, errors.New("no such url is in storage")
	}
	return url, nil
}

func (c *cacheRepo) GetByLong(longUrl string) (*domain.Url, error) {
	c.RLock()
	url, ok := (*c.mapLongKey)[longUrl]
	c.RUnlock()
	if !ok {
		return nil, errNoSuchUrl
	}
	return url, nil
}

func (c *cacheRepo) Get(url string) (*domain.Url, error) {
	URlptr, err := c.GetByShort(url)
	if err == nil {
		return URlptr, nil
	}
	URlptr, err = c.GetByLong(url)
	if err == nil {
		return URlptr, nil
	}
	return nil, err
}
