package services_test

import (
	"errors"
	"grovia/internal/models"
	"grovia/internal/mocks"
	"grovia/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetParentByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockParentRepository)

	expected := &models.Parent{ID: 1, Name: "Budi", PhoneNumber: "08111"}
	mockRepo.On("GetParentByID", 1, 1).Return(expected, nil)

	svc := services.NewParentService(mockRepo)
	result, err := svc.GetParentByID(1, 1)

	assert.NoError(t, err)
	assert.Equal(t, "Budi", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetParentByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockParentRepository)

	mockRepo.On("GetParentByID", 99, 1).Return(nil, gorm.ErrRecordNotFound)

	svc := services.NewParentService(mockRepo)
	result, err := svc.GetParentByID(99, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
	mockRepo.AssertExpectations(t)
}

func TestGetAllParent_Success(t *testing.T) {
	mockRepo := new(mocks.MockParentRepository)

	parents := []models.Parent{
		{ID: 1, Name: "Budi"},
		{ID: 2, Name: "Sari"},
	}
	mockRepo.On("GetAllParent", 1, 10, 0, "").Return(parents, 2, nil)

	svc := services.NewParentService(mockRepo)
	results, meta, err := svc.GetAllParent(1, "", "1", "10")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, 2, meta.TotalData)
	mockRepo.AssertExpectations(t)
}

func TestGetAllParent_RepoError(t *testing.T) {
	mockRepo := new(mocks.MockParentRepository)

	mockRepo.On("GetAllParent", 1, 10, 0, "").Return([]models.Parent{}, 0, errors.New("db error"))

	svc := services.NewParentService(mockRepo)
	results, meta, err := svc.GetAllParent(1, "", "1", "10")

	assert.Nil(t, results)
	assert.Nil(t, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Gagal mengambil")
	mockRepo.AssertExpectations(t)
}

func TestCheckPhoneExists_Found(t *testing.T) {
	mockRepo := new(mocks.MockParentRepository)

	expected := &models.Parent{ID: 1, PhoneNumber: "08123456789"}
	mockRepo.On("FindParentByPhoneNumber", "08123456789").Return(expected, nil)

	svc := services.NewParentService(mockRepo)
	result, err := svc.CheckPhoneExists("08123456789")

	assert.NoError(t, err)
	assert.Equal(t, "08123456789", result.PhoneNumber)
	mockRepo.AssertExpectations(t)
}

func TestCheckPhoneExists_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockParentRepository)

	mockRepo.On("FindParentByPhoneNumber", "0000").Return(nil, gorm.ErrRecordNotFound)

	svc := services.NewParentService(mockRepo)
	result, err := svc.CheckPhoneExists("0000")

	assert.Nil(t, result)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteParentByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockParentRepository)

	mockRepo.On("DeleteParentByID", 1, 1, 10).Return(nil)

	svc := services.NewParentService(mockRepo)
	err := svc.DeleteParentByID(1, 1, 10)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteParentByID_Failed(t *testing.T) {
	mockRepo := new(mocks.MockParentRepository)

	mockRepo.On("DeleteParentByID", 99, 1, 10).Return(errors.New("not found"))

	svc := services.NewParentService(mockRepo)
	err := svc.DeleteParentByID(99, 1, 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Gagal menghapus")
	mockRepo.AssertExpectations(t)
}

func TestGetAllParentAllLocation_Success(t *testing.T) {
	mockRepo := new(mocks.MockParentRepository)

	parents := []models.Parent{{ID: 1, Name: "Budi"}}
	mockRepo.On("GetAllParentAllLocation", "", 10, 0).Return(parents, 1, nil)

	svc := services.NewParentService(mockRepo)
	results, meta, err := svc.GetAllParentAllLocation("", "1", "10")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, 1, meta.TotalData)
	mockRepo.AssertExpectations(t)
}