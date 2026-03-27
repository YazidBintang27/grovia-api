package services_test

import (
	"encoding/json"
	"errors"
	"grovia/internal/dto/requests"
	"grovia/internal/mocks"
	"grovia/internal/models"
	"grovia/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newMockMLServer(t *testing.T, statusCode int, body map[string]any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(body)
	}))
}

func TestGetPredictByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockPredictRepository)

	expected := &models.Predict{ID: 1, NutritionalStatus: "Normal"}
	mockRepo.On("GetPredictByID", 1).Return(expected, nil)

	svc := services.NewPredictService(mockRepo, "http://localhost")
	result, err := svc.GetPredictByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "Normal", result.NutritionalStatus)
	mockRepo.AssertExpectations(t)
}

func TestGetPredictByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockPredictRepository)

	mockRepo.On("GetPredictByID", 99).Return(nil, errors.New("not found"))

	svc := services.NewPredictService(mockRepo, "http://localhost")
	result, err := svc.GetPredictByID(99)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
	mockRepo.AssertExpectations(t)
}

func TestGetAllPredict_Success(t *testing.T) {
	mockRepo := new(mocks.MockPredictRepository)

	predicts := []models.Predict{{ID: 1}, {ID: 2}}
	mockRepo.On("GetAllPredict", 1, 10, 0).Return(predicts, 2, nil)

	svc := services.NewPredictService(mockRepo, "http://localhost")
	results, meta, err := svc.GetAllPredict(1, "1", "10")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, 2, meta.TotalData)
	mockRepo.AssertExpectations(t)
}

func TestDeletePredictByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockPredictRepository)

	mockRepo.On("DeletePredictByID", 1, 1, 10).Return(nil)

	svc := services.NewPredictService(mockRepo, "http://localhost")
	err := svc.DeletePredictByID(1, 1, 10)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateIndividualPredict_MLAPISuccess(t *testing.T) {
	mockRepo := new(mocks.MockPredictRepository)

	// Mock ML server
	mlServer := newMockMLServer(t, 200, map[string]any{
		"zscore":            -1.5,
		"nutritionalStatus": "Normal",
	})
	defer mlServer.Close()

	birthdate := time.Now().AddDate(-2, 0, 0)
	req := requests.CreateToddlerRequest{
		Name:              "Anak Test",
		Sex:               "male",
		Birthdate:         birthdate,
		Height:            85.0,
		NutritionalStatus: "normal",
		LocationID:        2,
		PhoneNumber:       "0877749284",
	}

	savedPredict := &models.Predict{
		ID:                1,
		NutritionalStatus: "Normal",
		Zscore:            -1.5,
	}
	mockRepo.On("CreateIndividualPredict", mock.AnythingOfType("*models.Predict"), 1, 5).
		Return(savedPredict, nil)

	svc := services.NewPredictService(mockRepo, mlServer.URL)
	result, err := svc.CreateIndividualPredict(req, 1, 5, 10)

	assert.NoError(t, err)
	assert.Equal(t, "Normal", result.NutritionalStatus)
	mockRepo.AssertExpectations(t)
}

func TestCreateIndividualPredict_MLAPIError(t *testing.T) {
	mockRepo := new(mocks.MockPredictRepository)

	// ML server return error
	mlServer := newMockMLServer(t, 500, map[string]any{})
	defer mlServer.Close()

	birthdate := time.Now().AddDate(-2, 0, 0)
	req := requests.CreateToddlerRequest{
		Name:              "Anak Test",
		Sex:               "male",
		Birthdate:         birthdate,
		Height:            85.0,
		NutritionalStatus: "normal",
		LocationID:        2,
		PhoneNumber:       "0877749284",
	}

	svc := services.NewPredictService(mockRepo, mlServer.URL)
	result, err := svc.CreateIndividualPredict(req, 1, 5, 10)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ML API error")
}
