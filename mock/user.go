package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/thepoly/shuttletracker"
)

// UserService implements a mock of shuttletracker.UserService.
type UserService struct {
	mock.Mock
}

// UserExists returns whether the User exists.
func (us *UserService) UserExists(username string) (bool, error) {
	args := us.Called(username)
	return args.Bool(0), args.Error(1)
}

// Users gets all Users.
func (us *UserService) Users() ([]*shuttletracker.User, error) {
	args := us.Called()
	return args.Get(0).([]*shuttletracker.User), args.Error(1)
}

// CreateUser creates a User.
func (us *UserService) CreateUser(user *shuttletracker.User) error {
	args := us.Called(user)
	return args.Error(0)
}

// DeleteUser deletes a User.
func (us *UserService) DeleteUser(username string) error {
	args := us.Called(username)
	return args.Error(0)
}
