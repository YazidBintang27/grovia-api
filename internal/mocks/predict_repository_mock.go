package mocks

import (
	"grovia/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockPredictRepository struct {
	mock.Mock
}

func (m *MockPredictRepository) CreateIndividualPredict(predict *models.Predict, locationID, toddlerID int) (*models.Predict, error) {
	args := m.Called(predict, locationID, toddlerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Predict), args.Error(1)
}

func (m *MockPredictRepository) GetAllPredict(locationID, limit, offset int) ([]models.Predict, int, error) {
	args := m.Called(locationID, limit, offset)
	return args.Get(0).([]models.Predict), args.Int(1), args.Error(2)
}

func (m *MockPredictRepository) GetAllPredictByToddlerID(locationID, toddlerID int) ([]models.Predict, error) {
	args := m.Called(locationID, toddlerID)
	return args.Get(0).([]models.Predict), args.Error(1)
}

func (m *MockPredictRepository) GetPredictByID(id int) (*models.Predict, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Predict), args.Error(1)
}

func (m *MockPredictRepository) UpdatePredictByID(id int, predict *models.Predict) (*models.Predict, error) {
	args := m.Called(id, predict)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Predict), args.Error(1)
}

func (m *MockPredictRepository) DeletePredictByID(id, locationID, userID int) error {
	args := m.Called(id, locationID, userID)
	return args.Error(0)
}

func (m *MockPredictRepository) GetAllPredictAllLocation(limit, offset int) ([]models.Predict, int, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.Predict), args.Int(1), args.Error(2)
}