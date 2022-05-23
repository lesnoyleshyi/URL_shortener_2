package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"runtime"
)

func New() (*pgxpool.Pool, error) {
	connstr := os.Getenv("PG_CONNSTR")
	connConf, err := pgxpool.ParseConfig(connstr)
	if err != nil {
		return nil, fmt.Errorf("unable configure postgres connection: %w", err)
	}

	//another params could be defined here
	connConf.MaxConns = int32(runtime.NumCPU())

	pool, err := pgxpool.ConnectConfig(context.TODO(), connConf)
	if err != nil {
		return nil, fmt.Errorf("unable connect postgres: %w", err)
	}
	return pool, nil
}
