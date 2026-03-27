package mocks

import (
	"grovia/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) FindByPhoneNumber(phoneNumber string) (*models.User, error) {
	args := m.Called(phoneNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthRepository) ResetPassword(phoneNumber, newPassword string) error {
	args := m.Called(phoneNumber, newPassword)
	return args.Error(0)
}

func (m *MockAuthRepository) FindByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
