package repository

import (
	"URL_shortener_2/pkg/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type dbRepo struct {
	*pgxpool.Pool
}

func newDbRepo() *dbRepo {
	pgpool, err := postgres.New()
	if err != nil {
		log.Fatal(err)
	}
	return &dbRepo{Pool: pgpool}
}

func (r dbRepo) Save(longUrl string) {

}

func (r dbRepo) Get(longUrl string) string {
	return ""
}
