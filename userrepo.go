package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

// Create a User to the storage.
func (ur UserRepository) Create(u *User) error {
	stmt, err := ur.db.Prepare("INSERT INTO users (uuid, name, created_at) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	u.UUID = NextUUID()
	u.CreatedAt = time.Now()

	err = stmt.QueryRow(u.UUID, u.Name, u.CreatedAt).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

// ByID get a User from its id.
func (ur UserRepository) ByID(id int) (*User, error) {
	stmt, err := ur.db.Prepare("SELECT id, uuid, name, created_at, updated_at, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	u, err := scanRowUser(stmt.QueryRow(id))
	if errors.Is(err, sql.ErrNoRows) {
		return &User{}, ErrUserNotFound
	}

	if err != nil {
		return &User{}, err
	}

	return u, nil
}

// countAll return total of Users in storage.
func (ur UserRepository) countAll(f *Filter) (int, error) {
	stmt, err := ur.db.Prepare("SELECT COUNT (*) FROM users WHERE deleted_at IS NULL")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var n int
	err = stmt.QueryRow().Scan(&n)
	if err != nil {
		return 0, err
	}

	return n, nil
}

// All get all Users filterd per page and limit.
func (ur UserRepository) All(f *Filter) (FilteredResults, error) {
	query := "SELECT id, uuid, name, created_at, updated_at, deleted_at FROM users WHERE deleted_at IS NULL"
	query += " " + fmt.Sprintf("ORDER BY %s %s", f.Sort, f.Direction)
	query += " " + limitOffset(f.Limit, f.Page)

	stmt, err := ur.db.Prepare(query)
	if err != nil {
		return FilteredResults{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return FilteredResults{}, err
	}
	defer rows.Close()

	users := make(Users, 0)

	for rows.Next() {
		u, err := scanRowUser(rows)
		if err != nil {
			return FilteredResults{}, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return FilteredResults{}, err
	}

	// Get total rows to calculate total pages.
	totalRows, err := ur.countAll(f)
	if err != nil {
		return FilteredResults{}, err
	}

	return f.Paginate(users, totalRows), nil
}

// Delete soft deleted by User ID.
func (ur UserRepository) Delete(id int) error {
	stmt, err := ur.db.Prepare("UPDATE users SET deleted_at = $1 WHERE id = $2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

// scanRowUser return nulled fields of the domain object User parsed.
func scanRowUser(s scanner) (*User, error) {
	var updatedAtNull, deletedAtNull sql.NullTime
	m := &User{}

	err := s.Scan(
		&m.ID,
		&m.UUID,
		&m.Name,
		&m.CreatedAt,
		&updatedAtNull,
		&deletedAtNull,
	)
	if err != nil {
		return &User{}, err
	}

	m.UpdatedAt = updatedAtNull.Time
	m.DeletedAt = deletedAtNull.Time

	return m, nil
}

// timeToNull helper to try empty time fields.
func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}

	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}

// limitOffset returns a SQL string for a given limit & offset.
func limitOffset(limit, page int) string {
	if limit == 0 && page == 0 {
		return ""
	}

	offset := page*limit - limit
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}
