//go:build integration
// +build integration

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
)

func TestCreateUser(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	db := openDB(t)
	defer closeDB(t, db)

	ur := UserRepository{db: db}

	input := &User{
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

	db := openDB(t)
	defer closeDB(t, db)
	insertUsersData(t, db)

	ur := UserRepository{db: db}
	userID := 1

	if err := ur.Delete(userID); err != nil {
		t.Fatal(err)
	}

	got, err := onlyTrashedByID(db, userID)
	if err != nil {
		t.Fatal(err)
	}

	if got.DeletedAt.IsZero() {
		t.Error("expected deleted at")
	}
}

func onlyTrashedByID(db *sql.DB, id int) (*User, error) {
	stmt, err := db.Prepare("SELECT * FROM users WHERE id = $1 AND deleted_at IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	u, err := scanRowUser(stmt.QueryRow(id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}

func truncate(db *sql.DB, table string) error {
	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", table)
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}

	return nil
}

func cleanUsersData(t *testing.T) {
	db := openDB(t)
	defer closeDB(t, db)

	err := truncate(db, "users")
	if err != nil {
		t.Fatal(err)
	}
}

func insertUsersData(t *testing.T, db *sql.DB) {
	ur := UserRepository{db: db}

	if err := ur.Create(&User{
		Name: "John",
	}); err != nil {
		t.Fatal(err)
	}

	if err := ur.Create(&User{
		Name: "Jane",
	}); err != nil {
		t.Fatal(err)
	}
}
