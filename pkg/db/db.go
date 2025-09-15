package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPool(dsn string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool - %w", err)
	}

	return dbPool, nil
}
