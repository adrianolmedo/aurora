package main

import "fmt"

// Storage represents all repositories.
type Storage struct {
	UserRepo UserRepository
}

// NewStorage init postgres database conecton, build the Storage and return
// it its pointer.
func NewStorage(cfg Config) (*Storage, error) {
	if cfg.EngineDB == "" {
		return nil, fmt.Errorf("database engine not selected")
	}

	if cfg.EngineDB == "postgres" {
		db, err := newDB(cfg)
		if err != nil {
			return nil, fmt.Errorf("postgres: %v", err)
		}

		return &Storage{
			UserRepo: UserRepository{db: db},
		}, nil
	}

	return nil, fmt.Errorf("database engine '%s' not implemented", cfg.EngineDB)
}
