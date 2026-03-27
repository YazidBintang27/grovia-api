package mocks

import (
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"

	"github.com/stretchr/testify/mock"
)

type MockPredictService struct {
	mock.Mock
}

func (m *MockPredictService) CreateIndividualPredict(req requests.CreateToddlerRequest, locationID, toddlerID, userID int) (*responses.PredictResponse, error) {
	args := m.Called(req, locationID, toddlerID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.PredictResponse), args.Error(1)
}

func (m *MockPredictService) CreateGroupPredict(filePath string) ([]byte, error) {
	args := m.Called(filePath)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockPredictService) GetAllPredict(locationID int, pageStr, limitStr string) ([]responses.PredictResponse, *responses.PaginationMeta, error) {
	args := m.Called(locationID, pageStr, limitStr)
	return args.Get(0).([]responses.PredictResponse), args.Get(1).(*responses.PaginationMeta), args.Error(2)
}

func (m *MockPredictService) GetAllPredictByToddlerID(locationID, toddlerID int) ([]responses.PredictResponse, error) {
	args := m.Called(locationID, toddlerID)
	return args.Get(0).([]responses.PredictResponse), args.Error(1)
}

func (m *MockPredictService) GetPredictByID(id int) (*responses.PredictResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.PredictResponse), args.Error(1)
}

func (m *MockPredictService) UpdatePredictByID(id int, req *requests.UpdatePredictRequest) (*responses.PredictResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.PredictResponse), args.Error(1)
}

func (m *MockPredictService) DeletePredictByID(id, locationID, userID int) error {
	args := m.Called(id, locationID, userID)
	return args.Error(0)
}

func (m *MockPredictService) GetAllPredictAllLocation(pageStr, limitStr string) ([]responses.PredictResponse, *responses.PaginationMeta, error) {
	args := m.Called(pageStr, limitStr)
	return args.Get(0).([]responses.PredictResponse), args.Get(1).(*responses.PaginationMeta), args.Error(2)
}