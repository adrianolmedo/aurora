//go:build integration
// +build integration

package main

import (
	"database/sql"
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

func TestUpdateUser(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersData(t)
	})

	db := openDB(t)
	defer closeDB(t, db)
	insertUsersData(t, db)

	input := User{
		ID:   1,
		Name: "Adrián",
	}

	ur := UserRepository{db: db}

	if err := ur.Update(input); err != nil {
		t.Fatal(err)
	}

	got, err := ur.ByID(input.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.Name != input.Name {
		t.Errorf("Name: want %s, got %s", input.Name, got.Name)
	}

	if got.CreatedAt.IsZero() {
		t.Error("expected created at")
	}

	if got.UpdatedAt.IsZero() {
		t.Error("expected updated at")
	}

	if !got.DeletedAt.IsZero() {
		t.Error("unexpected deleted at")
	}
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
