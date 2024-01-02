//go:build integration
// +build integration

package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	domain "github.com/adrianolmedo/aurora"

	"github.com/jackc/pgx/v5"
)

func TestCreateUser(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	conn := openDB(t)
	defer closeDB(t, conn)

	ur := UserRepo{conn: conn}

	input := &domain.User{
		Name: "Adrián",
	}

	if err := ur.Create(input); err != nil {
		t.Fatal(err)
	}

	got, err := ur.ByID(1)
	if err != nil {
		t.Fatal(err)
	}

	if got.CreatedAt.IsZero() {
		t.Error("expected created at")
	}

	if !got.UpdatedAt.IsZero() {
		t.Error("unexpected updated at")
	}

	if !got.DeletedAt.IsZero() {
		t.Error("unexpected deleted at")
	}
}

func TestDeleteUser(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	conn := openDB(t)
	defer closeDB(t, conn)
	insertUsersData(t, conn)

	ur := UserRepo{conn: conn}
	userID := 1

	if err := ur.Delete(userID); err != nil {
		t.Fatal(err)
	}

	got, err := onlyTrashedByID(conn, userID)
	if err != nil {
		t.Fatal(err)
	}

	if got.DeletedAt.IsZero() {
		t.Error("expected deleted at")
	}
}

func onlyTrashedByID(conn *pgx.Conn, id int) (*domain.User, error) {
	var updatedAtNull, deletedAtNull sql.NullTime
	m := &domain.User{}

	err := conn.QueryRow(context.Background(), "SELECT id, uuid, name, created_at, updated_at, deleted_at FROM users WHERE id = $1 AND deleted_at IS NOT NULL", id).
		Scan(&m.ID, &m.UUID, &m.Name, &m.CreatedAt, &updatedAtNull, &deletedAtNull)
	if err != nil {
		return nil, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	m.UpdatedAt = updatedAtNull.Time
	m.DeletedAt = deletedAtNull.Time

	return m, nil
}

func truncate(conn *pgx.Conn, table string) error {
	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", table)
	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}

	return nil
}

func cleanUsersData(t *testing.T) {
	conn := openDB(t)
	defer closeDB(t, conn)

	err := truncate(conn, "users")
	if err != nil {
		t.Fatal(err)
	}
}

func insertUsersData(t *testing.T, conn *pgx.Conn) {
	ur := UserRepo{conn: conn}

	if err := ur.Create(&domain.User{
		Name: "John",
	}); err != nil {
		t.Fatal(err)
	}

	if err := ur.Create(&domain.User{
		Name: "Jane",
	}); err != nil {
		t.Fatal(err)
	}
}
