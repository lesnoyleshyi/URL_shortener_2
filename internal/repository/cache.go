package repository

import "sync"

type cacheRepo struct {
	sync.RWMutex
	storage *map[string]string
}

func newCache() *cacheRepo {
	m := make(map[string]string)
	cacheStruct := cacheRepo{
		storage: &m,
	}
	return &cacheStruct
}

func (c cacheRepo) Save(longUrl string) {

}

func (c cacheRepo) Get(shortUrl string) string {
	return ""
}
