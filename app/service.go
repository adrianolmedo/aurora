package app

import (
	domain "github.com/adrianolmedo/aurora"
	"github.com/adrianolmedo/aurora/storage"
)

type Service struct {
	User UserService
}

func NewService(storage *storage.Storage) *Service {
	return &Service{
		User: UserService{repo: storage.UserRepo},
	}
}

type UserService struct {
	repo storage.UserRepo
}

// SignUp to register a User.
func (us UserService) SignUp(u *domain.User) error {
	err := signUp(u)
	if err != nil {
		return err
	}

	return us.repo.Create(u)
}

// signUp applicaction logic for regitser a User. Has been split into
// a smaller function for unit testing purposes, and it should do so for
// the other methods of the Service.
func signUp(u *domain.User) error {
	err := u.Validate()
	if err != nil {
		return err
	}

	return nil
}

// List get list of Users.
func (usp UserService) List(f *domain.Filter) (domain.FilteredResults, error) {
	return usp.repo.All(f)
}

// Find a User by its ID.
func (us UserService) Find(id int) (*domain.User, error) {
	if id == 0 {
		return &domain.User{}, domain.ErrUserNotFound
	}

	return us.repo.ByID(id)
}

// Remove mark User as deleted by its ID.
func (us UserService) Remove(id int) error {
	if id == 0 {
		return domain.ErrUserNotFound
	}

	return us.repo.Delete(id)
}
