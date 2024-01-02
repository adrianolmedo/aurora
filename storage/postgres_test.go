//go:build integration
// +build integration

package storage

import (
	"context"
	"flag"
	"testing"

	"github.com/adrianolmedo/aurora"

	"github.com/jackc/pgx/v5"
)

// $ go test -v -tags integration -args -dbhost 127.0.0.1 -dbport 5432 -dbuser username -dbname foodb -dbpass 12345
var (
	dbhost = flag.String("dbhost", "", "Database host.")
	dbport = flag.String("dbport", "", "Database port.")
	dbuser = flag.String("dbuser", "", "Database user.")
	dbpass = flag.String("dbpass", "", "Database password.")
	dbname = flag.String("dbname", "", "Database name.")
)

// TestDB test for open & close database.
func TestDB(t *testing.T) {
	conn := openDB(t)
	closeDB(t, conn)
}

func openDB(t *testing.T) *pgx.Conn {
	dbcfg := aurora.Config{
		HostDB:     *dbhost,
		PortDB:     *dbport,
		UserDB:     *dbuser,
		PasswordDB: *dbpass,
		NameDB:     *dbname,
	}

	conn, err := newDB(dbcfg)
	if err != nil {
		t.Fatal(err)
	}

	return conn
}

func closeDB(t *testing.T, conn *pgx.Conn) {
	err := conn.Close(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
