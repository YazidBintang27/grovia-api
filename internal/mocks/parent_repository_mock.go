package mocks

import (
	"grovia/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockParentRepository struct {
	mock.Mock
}

func (m *MockParentRepository) CreateParent(parent *models.Parent) (*models.Parent, error) {
	args := m.Called(parent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Parent), args.Error(1)
}

func (m *MockParentRepository) GetAllParent(locationID, limit, offset int, name string) ([]models.Parent, int, error) {
	args := m.Called(locationID, limit, offset, name)
	return args.Get(0).([]models.Parent), args.Int(1), args.Error(2)
}

func (m *MockParentRepository) GetParentByID(id, locationID int) (*models.Parent, error) {
	args := m.Called(id, locationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Parent), args.Error(1)
}

func (m *MockParentRepository) UpdateParentByID(id, locationID int, parent *models.Parent) (*models.Parent, error) {
	args := m.Called(id, locationID, parent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Parent), args.Error(1)
}

func (m *MockParentRepository) DeleteParentByID(id, locationID, userID int) error {
	args := m.Called(id, locationID, userID)
	return args.Error(0)
}

func (m *MockParentRepository) FindParentByPhoneNumber(phoneNumber string) (*models.Parent, error) {
	args := m.Called(phoneNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Parent), args.Error(1)
}

func (m *MockParentRepository) GetAllParentAllLocation(name string, limit, offset int) ([]models.Parent, int, error) {
	args := m.Called(name, limit, offset)
	return args.Get(0).([]models.Parent), args.Int(1), args.Error(2)
}