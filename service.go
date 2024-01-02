package main

type Service struct {
	User UserService
}

func NewService(storage *Storage) *Service {
	return &Service{
		User: UserService{repo: storage.UserRepo},
	}
}

type UserService struct {
	repo UserRepo
}

// SignUp to register a User.
func (us UserService) SignUp(u *User) error {
	err := signUp(u)
	if err != nil {
		return err
	}

	return us.repo.Create(u)
}

// signUp applicaction logic for regitser a User. Has been split into
// a smaller function for unit testing purposes, and it should do so for
// the other methods of the Service.
func signUp(u *User) error {
	err := u.Validate()
	if err != nil {
		return err
	}

	return nil
}

// List get list of Users.
func (usp UserService) List(f *Filter) (FilteredResults, error) {
	return usp.repo.All(f)
}

// Find a User by its ID.
func (us UserService) Find(id int) (*User, error) {
	if id == 0 {
		return &User{}, ErrUserNotFound
	}

	return us.repo.ByID(id)
}

// Remove mark User as deleted by its ID.
func (us UserService) Remove(id int) error {
	if id == 0 {
		return ErrUserNotFound
	}

	return us.repo.Delete(id)
}
