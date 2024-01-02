package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func newDB(cfg Config) (*pgx.Conn, error) {
	// postgres://user:password@host:port/dbname?sslmode=disable
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.UserDB, cfg.PasswordDB, cfg.HostDB, cfg.PortDB, cfg.NameDB)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't do ping %v", err)
	}

	//defer conn.Close(context.Background())

	return conn, nil
}
