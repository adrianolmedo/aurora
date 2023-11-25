//go:build pgx
// +build pgx

package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestCreateUserPGX(t *testing.T) {
	t.Cleanup(func() {
		cleanUsersDataPGX(t)
	})

	conn := openPGX(t)
	defer closePGX(t, conn)

	ur := UserRepoPGX{conn: conn}

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

func truncatePGX(conn *pgx.Conn, table string) error {
	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", table)
	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}

	return nil
}

func cleanUsersDataPGX(t *testing.T) {
	conn := openPGX(t)
	defer closePGX(t, conn)

	err := truncatePGX(conn, "users")
	if err != nil {
		t.Fatal(err)
	}
}
