package mocks

import (
	"grovia/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUser(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetAllUser() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(id int, user *models.User) (*models.User, error) {
	args := m.Called(id, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) FindRoleById(id int) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) FindUsersByRole(roles []string, name string, locationID, limit, offset int) ([]models.User, int, error) {
	args := m.Called(roles, name, locationID, limit, offset)
	return args.Get(0).([]models.User), args.Int(1), args.Error(2)
}