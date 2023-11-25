//go:build pgx
// +build pgx

package main

import (
	"context"
	"flag"
	"testing"

	"github.com/jackc/pgx/v5"
)

// $ go test -v -tags pgx -args -dbhost 127.0.0.1 -dbport 5432 -dbuser username -dbname foodb -dbpass 12345
var (
	dbhost   = flag.String("dbhost", "", "Database host.")
	dbengine = flag.String("dbengine", "postgres", "Database engine, choose mysql or postgres.")
	dbport   = flag.String("dbport", "", "Database port.")
	dbuser   = flag.String("dbuser", "", "Database user.")
	dbpass   = flag.String("dbpass", "", "Database password.")
	dbname   = flag.String("dbname", "", "Database name.")
)

// TestDB test for open & close database.
func TestPGX(t *testing.T) {
	conn := openPGX(t)
	closePGX(t, conn)
}

func openPGX(t *testing.T) *pgx.Conn {
	dbcfg := Config{
		EngineDB:   *dbengine,
		HostDB:     *dbhost,
		PortDB:     *dbport,
		UserDB:     *dbuser,
		PasswordDB: *dbpass,
		NameDB:     *dbname,
	}

	conn, err := newDBWithPGX(dbcfg)
	if err != nil {
		t.Fatal(err)
	}

	return conn
}

func closePGX(t *testing.T, conn *pgx.Conn) {
	err := conn.Close(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
