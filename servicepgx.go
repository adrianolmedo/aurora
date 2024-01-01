package main

type ServicePGX struct {
	User UserServicePGX
}

func NewServicePGX(storage *StoragePGX) *ServicePGX {
	return &ServicePGX{
		User: UserServicePGX{repo: storage.UserRepo},
	}
}

type UserServicePGX struct {
	repo UserRepoPGX
}

// List get list of Users.
func (usp UserServicePGX) List(f *Filter) (FilteredResults, error) {
	return usp.repo.All(f)
}
