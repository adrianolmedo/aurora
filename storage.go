package main

import "fmt"

// Storage represents all repositories.
type Storage struct {
	UserRepo UserRepository
}

// StoragePGX represents all repositories.
type StoragePGX struct {
	UserRepo UserRepoPGX
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

// NewPGX return pointer of StoragePGX.
func NewPGX(cfg Config) (*StoragePGX, error) {
	if cfg.EngineDB == "" {
		return nil, fmt.Errorf("database engine not selected")
	}

	if cfg.EngineDB == "postgres" {
		conn, err := newDBWithPGX(cfg)
		if err != nil {
			return nil, fmt.Errorf("postgres: %v", err)
		}

		return &StoragePGX{
			UserRepo: UserRepoPGX{conn: conn},
		}, nil
	}

	return nil, fmt.Errorf("database engine '%s' not implemented", cfg.EngineDB)
}
