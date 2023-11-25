package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserRepoPGX struct {
	conn *pgx.Conn
}

func (r UserRepoPGX) Create(u *User) error {
	u.UUID = NextUUID()
	u.CreatedAt = time.Now()

	err := r.conn.QueryRow(context.Background(), "INSERT INTO users (uuid, name, created_at) VALUES ($1, $2, $3) RETURNING id", u.UUID, u.Name, u.CreatedAt).Scan(&u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r UserRepoPGX) ByID(id int) (*User, error) {
	var updatedAtNull, deletedAtNull sql.NullTime

	m := &User{}

	err := r.conn.QueryRow(context.Background(), "SELECT id, uuid, name, created_at, updated_at, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL", id).Scan(&m.ID, &m.UUID, &m.Name, &m.CreatedAt, &updatedAtNull, &deletedAtNull)
	if err != nil {
		return nil, err
	}
	return m, nil
}
