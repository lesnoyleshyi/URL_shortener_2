package repository

type Repository interface {
	Save(longUrl string)
	Get(shortUrl string) string
}

func New(storageType string) Repository {
	if storageType == "cache" {
		return newCache()
	} else {
		return newDbRepo()
	}
}
