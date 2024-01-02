package main

import "fmt"

// Storage represents all repositories.
type Storage struct {
	UserRepo UserRepo
}

// NewPGX return pointer of StoragePGX.
func NewStorage(cfg Config) (*Storage, error) {
	if cfg.EngineDB == "" {
		return nil, fmt.Errorf("database engine not selected")
	}

	if cfg.EngineDB == "postgres" {
		conn, err := newDB(cfg)
		if err != nil {
			return nil, fmt.Errorf("postgres: %v", err)
		}

		return &Storage{
			UserRepo: UserRepo{conn: conn},
		}, nil
	}

	return nil, fmt.Errorf("database engine '%s' not implemented", cfg.EngineDB)
}
