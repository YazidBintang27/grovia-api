package services_test

import (
	"context"
	"errors"
	"grovia/internal/dto/requests"
	"grovia/internal/models"
	"grovia/internal/mocks"
	"grovia/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetLocationByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)

	expected := &models.Location{ID: 1, Name: "Posyandu Mawar"}
	mockRepo.On("GetLocationByID", 1).Return(expected, nil)

	svc := services.NewLocationService(mockRepo, nil)
	result, err := svc.GetLocationByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "Posyandu Mawar", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetLocationByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)

	mockRepo.On("GetLocationByID", 99).Return(nil, gorm.ErrRecordNotFound)

	svc := services.NewLocationService(mockRepo, nil)
	result, err := svc.GetLocationByID(99)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
	mockRepo.AssertExpectations(t)
}

func TestGetAllLocation_Success(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)

	locations := []models.Location{
		{ID: 1, Name: "Posyandu A"},
		{ID: 2, Name: "Posyandu B"},
	}
	mockRepo.On("GetAllLocation", "", 10, 0).Return(locations, 2, nil)

	svc := services.NewLocationService(mockRepo, nil)
	results, meta, err := svc.GetAllLocation("", "1", "10")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, 2, meta.TotalData)
	assert.Equal(t, 1, meta.TotalPage)
	mockRepo.AssertExpectations(t)
}

func TestGetAllLocation_DefaultPagination(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)

	mockRepo.On("GetAllLocation", "", 1, 0).Return([]models.Location{}, 0, nil)

	svc := services.NewLocationService(mockRepo, nil)

	// page=0 dan limit=0 harus default ke 1
	results, meta, err := svc.GetAllLocation("", "0", "0")

	assert.NoError(t, err)
	assert.Empty(t, results)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 1, meta.Limit)
	mockRepo.AssertExpectations(t)
}

func TestDeleteLocationByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)

	mockRepo.On("DeleteLocationByID", 1, 1).Return(nil)

	svc := services.NewLocationService(mockRepo, nil)
	err := svc.DeleteLocationByID(1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteLocationByID_Failed(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)

	mockRepo.On("DeleteLocationByID", 99, 1).Return(errors.New("not found"))

	svc := services.NewLocationService(mockRepo, nil)
	err := svc.DeleteLocationByID(99, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Gagal menghapus")
	mockRepo.AssertExpectations(t)
}

func TestCreateLocation_InvalidRequest(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)
	svc := services.NewLocationService(mockRepo, nil)

	// Kirim request kosong — validasi harus gagal sebelum ke repo
	req := requests.LocationRequest{}
	result, err := svc.CreateLocation(context.Background(), req, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestUpdateLocationByID_InvalidRequest(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)
	svc := services.NewLocationService(mockRepo, nil)

	req := requests.LocationRequest{}
	result, err := svc.UpdateLocationByID(context.Background(), 1, 1, req)

	assert.Nil(t, result)
	assert.Error(t, err)
}