package storage

import (
	"fmt"

	"github.com/adrianolmedo/aurora"
)

// Storage represents all repositories.
type Storage struct {
	UserRepo UserRepo
}

// New return pointer of Storage.
func New(cfg aurora.Config) (*Storage, error) {
	conn, err := newDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: %v", err)
	}

	return &Storage{
		UserRepo: UserRepo{conn: conn},
	}, nil
}
