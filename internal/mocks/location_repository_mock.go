package mocks

import (
	"grovia/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockLocationRepository struct {
	mock.Mock
}

func (m *MockLocationRepository) CreateLocation(location *models.Location) (*models.Location, error) {
	args := m.Called(location)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Location), args.Error(1)
}

func (m *MockLocationRepository) GetAllLocation(name string, limit, offset int) ([]models.Location, int, error) {
	args := m.Called(name, limit, offset)
	return args.Get(0).([]models.Location), args.Int(1), args.Error(2)
}

func (m *MockLocationRepository) GetLocationByID(id int) (*models.Location, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Location), args.Error(1)
}

func (m *MockLocationRepository) UpdateLocationByID(id int, location *models.Location) (*models.Location, error) {
	args := m.Called(id, location)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Location), args.Error(1)
}

func (m *MockLocationRepository) DeleteLocationByID(id, userID int) error {
	args := m.Called(id, userID)
	return args.Error(0)
}