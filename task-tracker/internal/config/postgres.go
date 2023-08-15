package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

const pgDSN = "postgres://postgres:password@localhost:5435/task-tracker?sslmode=disable"

func ConnectToPostgres(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, pgDSN)
	if err != nil {
		return nil, fmt.Errorf("cant connect to postgres %w", err)
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("cant setup connection to postgres %w", err)
	}

	return conn, nil
}
