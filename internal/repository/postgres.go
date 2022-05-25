package repository

import (
	"URL_shortener_2/internal/domain"
	"URL_shortener_2/pkg/postgres"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
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

const timeoutSec = 1
const saveUrlQuery = `INSERT INTO urls (short_url, long_url, created_at)
								VALUES ($1, $2, $3);`
const getUrlQuery = `SELECT short_url, long_url, created_at FROM urls
						WHERE short_url = $1 OR long_url = $1;`

func (r dbRepo) Save(url *domain.Url) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeoutSec)
	defer cancel()
	txOpts := pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	}

	tx, err := r.BeginTx(ctx, txOpts)
	if err != nil {
		return fmt.Errorf("unable begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, err = r.Exec(ctx, saveUrlQuery, url.ShortUrl, url.LongUrl, url.CreatedAt)
	if err != nil {
		return fmt.Errorf("unable execute sql query: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("unable commit transaction: %w", err)
	}
	return nil
}

func (r dbRepo) Get(url string) (*domain.Url, error) {
	urlStruct := domain.Url{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeoutSec)
	defer cancel()
	txOpts := pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadOnly,
		DeferrableMode: pgx.NotDeferrable,
	}

	tx, err := r.BeginTx(ctx, txOpts)
	if err != nil {
		return nil, fmt.Errorf("unable begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	row := tx.QueryRow(ctx, getUrlQuery, url)
	err = row.Scan(&urlStruct.ShortUrl, &urlStruct.LongUrl, &urlStruct.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("unable retrieve data from database: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("unable commit transaction: %w", err)
	}
	return &urlStruct, nil
}
