package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectToPostgres(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, os.Getenv(`PG_DSN`))
	if err != nil {
		return nil, fmt.Errorf("cant connect to postgres %w", err)
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("cant setup connection to postgres %w", err)
	}

	return conn, nil
}
