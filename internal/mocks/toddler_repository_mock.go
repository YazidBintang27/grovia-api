package mocks

import (
	"grovia/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockToddlerRepository struct {
	mock.Mock
}

func (m *MockToddlerRepository) CreateToddler(toddler *models.Toddler) (*models.Toddler, error) {
	args := m.Called(toddler)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Toddler), args.Error(1)
}

func (m *MockToddlerRepository) GetAllToddler(locationID, limit, offset int, name string) ([]models.Toddler, int, error) {
	args := m.Called(locationID, limit, offset, name)
	return args.Get(0).([]models.Toddler), args.Int(1), args.Error(2)
}

func (m *MockToddlerRepository) GetToddlerByID(id, locationID int) (*models.Toddler, error) {
	args := m.Called(id, locationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Toddler), args.Error(1)
}

func (m *MockToddlerRepository) UpdateToddlerByID(id, locationID int, toddler *models.Toddler) (*models.Toddler, error) {
	args := m.Called(id, locationID, toddler)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Toddler), args.Error(1)
}

func (m *MockToddlerRepository) DeleteToddlerByID(id, locationID, userID int) error {
	args := m.Called(id, locationID, userID)
	return args.Error(0)
}

func (m *MockToddlerRepository) FindToddlerByName(parentID int, name string) (bool, *models.Toddler, error) {
	args := m.Called(parentID, name)
	if args.Get(1) == nil {
		return args.Bool(0), nil, args.Error(2)
	}
	return args.Bool(0), args.Get(1).(*models.Toddler), args.Error(2)
}

func (m *MockToddlerRepository) GetAllToddlerAllLocation(name string, limit, offset int) ([]models.Toddler, int, error) {
	args := m.Called(name, limit, offset)
	return args.Get(0).([]models.Toddler), args.Int(1), args.Error(2)
}