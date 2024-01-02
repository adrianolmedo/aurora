package aurora

import (
	"errors"
	"time"

	"github.com/pborman/uuid"
)

var ErrUserNotFound = errors.New("user not found")

// User model.
type User struct {
	ID   int    `json:"id"`
	UUID string `json:"uuid"`
	Name string `json:"name"`

	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

// Validate return error if certain fields there are empty.
func (u User) Validate() error {
	if u.Name == "" {
		return errors.New("the name can't be empty")
	}

	return nil
}

// Users a collection of User.
type Users []*User

// IsEmpty return true if is empty.
func (us Users) IsEmpty() bool {
	return len(us) == 0
}

// NextUUID generates a new UUID.
func NextUUID() string {
	return uuid.New()
}
