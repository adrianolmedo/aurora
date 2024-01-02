package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	conn *pgx.Conn
}

func (r UserRepo) Create(u *User) error {
	u.UUID = NextUUID()
	u.CreatedAt = time.Now()

	err := r.conn.QueryRow(context.Background(), "INSERT INTO users (uuid, name, created_at) VALUES ($1, $2, $3) RETURNING id", u.UUID, u.Name, u.CreatedAt).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepo) ByID(id int) (*User, error) {
	var updatedAtNull, deletedAtNull sql.NullTime

	m := &User{}

	err := r.conn.QueryRow(context.Background(), "SELECT id, uuid, name, created_at, updated_at, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL", id).
		Scan(&m.ID, &m.UUID, &m.Name, &m.CreatedAt, &updatedAtNull, &deletedAtNull)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	m.UpdatedAt = updatedAtNull.Time
	m.DeletedAt = deletedAtNull.Time

	return m, nil
}

func (r UserRepo) All(f *Filter) (FilteredResults, error) {
	query := "SELECT id, uuid, name, created_at, updated_at, deleted_at FROM users WHERE deleted_at IS NULL"
	query += " " + fmt.Sprintf("ORDER BY %s %s", f.Sort, f.Direction)
	query += " " + limitOffset(f.Limit, f.Page)

	rows, err := r.conn.Query(context.Background(), query)
	if err != nil {
		return FilteredResults{}, err
	}

	users := make(Users, 0)

	for rows.Next() {
		var updatedAtNull, deletedAtNull sql.NullTime
		m := &User{}

		err := rows.Scan(
			&m.ID,
			&m.UUID,
			&m.Name,
			&m.CreatedAt,
			&updatedAtNull,
			&deletedAtNull,
		)
		if err != nil {
			return FilteredResults{}, err
		}

		m.UpdatedAt = updatedAtNull.Time
		m.DeletedAt = deletedAtNull.Time

		users = append(users, m)
	}

	if err := rows.Err(); err != nil {
		return FilteredResults{}, err
	}

	// Get total rows to calculate total pages.
	totalRows, err := r.countAll(f)
	if err != nil {
		return FilteredResults{}, err
	}

	return f.Paginate(users, totalRows), nil

}

// countAll return total of Users in storage.
func (r UserRepo) countAll(f *Filter) (int, error) {
	var n int

	err := r.conn.QueryRow(context.Background(), "SELECT COUNT (*) FROM users WHERE deleted_at IS NULL").Scan(&n)
	if err != nil {
		return 0, err
	}

	return n, nil
}

// limitOffset returns a SQL string for a given limit & offset.
func limitOffset(limit, page int) string {
	if limit == 0 && page == 0 {
		return ""
	}

	offset := page*limit - limit
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

func (r UserRepo) Delete(id int) error {
	result, err := r.conn.Exec(context.Background(), "UPDATE users SET deleted_at = $1 WHERE id = $2", time.Now(), id)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}
