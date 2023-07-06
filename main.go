package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v3"
)

func main() {
	fs := flag.NewFlagSet("aurora", flag.ExitOnError)
	var (
		port     = fs.String("port", "80", "Internal container port.")
		dbhost   = fs.String("dbhost", "", "Database host.")
		dbengine = fs.String("dbengine", "postgres", "Database engine, choose mysql or postgres.")
		dbport   = fs.String("dbport", "", "Database port.")
		dbuser   = fs.String("dbuser", "", "Database user.")
		dbpass   = fs.String("dbpass", "", "Database password.")
		dbname   = fs.String("dbname", "", "Database name.")
	)

	// Pass env vars to flags.
	ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())

	err := run(Config{
		Port:       *port,
		EngineDB:   *dbengine,
		HostDB:     *dbhost,
		PortDB:     *dbport,
		UserDB:     *dbuser,
		PasswordDB: *dbpass,
		NameDB:     *dbname,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(cfg Config) error {
	strg, err := NewStorage(cfg)
	if err != nil {
		return fmt.Errorf("error from storage: %v", err)
	}

	return router(strg).Listen(":" + cfg.Port)
}
