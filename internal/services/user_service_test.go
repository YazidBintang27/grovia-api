package services_test

import (
	"errors"
	"grovia/internal/models"
	"grovia/internal/mocks"
	"grovia/internal/services"
	"grovia/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	expected := &models.User{ID: 1, Name: "Admin"}
	mockRepo.On("GetUser", 1).Return(expected, nil)

	svc := services.NewUserService(mockRepo, nil)
	result, err := svc.GetCurrentUser(1)

	assert.NoError(t, err)
	assert.Equal(t, "Admin", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetCurrentUser_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	mockRepo.On("GetUser", 99).Return(nil, errors.New("not found"))

	svc := services.NewUserService(mockRepo, nil)
	result, err := svc.GetCurrentUser(99)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
	mockRepo.AssertExpectations(t)
}

func TestDeleteCurrentUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	mockRepo.On("DeleteUser", 1).Return(nil)

	svc := services.NewUserService(mockRepo, nil)
	err := svc.DeleteCurrentUser(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteCurrentUser_Failed(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	mockRepo.On("DeleteUser", 99).Return(errors.New("not found"))

	svc := services.NewUserService(mockRepo, nil)
	err := svc.DeleteCurrentUser(99)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Gagal menghapus")
	mockRepo.AssertExpectations(t)
}

func TestGetUsersByRole_AdminSuccess(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	users := []models.User{
		{ID: 1, Role: pkg.RoleKader},
		{ID: 2, Role: pkg.RoleKepalaPosyandu},
	}
	expectedRoles := []string{pkg.RoleKepalaPosyandu, pkg.RoleKader}
	mockRepo.On("FindUsersByRole", expectedRoles, "", 1, 10, 0).Return(users, 2, nil)

	svc := services.NewUserService(mockRepo, nil)
	results, meta, err := svc.GetUsersByRole(pkg.RoleAdmin, "", "1", "10", 1)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, 2, meta.TotalData)
	mockRepo.AssertExpectations(t)
}

func TestGetUsersByRole_ForbiddenRole(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := services.NewUserService(mockRepo, nil)

	results, meta, err := svc.GetUsersByRole(pkg.RoleKader, "", "1", "10", 1)

	assert.Nil(t, results)
	assert.Nil(t, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak diizinkan")
}

func TestDeleteUserByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	mockRepo.On("FindRoleById", 2).Return(pkg.RoleKader, nil)
	mockRepo.On("DeleteUser", 2).Return(nil)

	svc := services.NewUserService(mockRepo, nil)
	err := svc.DeleteUserByID(2, pkg.RoleAdmin)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUserByID_Forbidden(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	// Kader mencoba hapus kepala posyandu — tidak boleh
	mockRepo.On("FindRoleById", 2).Return(pkg.RoleKepalaPosyandu, nil)

	svc := services.NewUserService(mockRepo, nil)
	err := svc.DeleteUserByID(2, pkg.RoleKader)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Tidak memiliki akses")
	mockRepo.AssertExpectations(t)
}

func TestGetUserById_ForbiddenAccess(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	mockRepo.On("FindRoleById", 5).Return(pkg.RoleAdmin, nil)

	svc := services.NewUserService(mockRepo, nil)
	result, err := svc.GetUserById(5, pkg.RoleKader)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Tidak memiliki akses")
	mockRepo.AssertExpectations(t)
}

func TestGetUserById_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	mockRepo.On("FindRoleById", 99).Return("", errors.New("not found"))

	svc := services.NewUserService(mockRepo, nil)
	result, err := svc.GetUserById(99, pkg.RoleAdmin)

	assert.Nil(t, result)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}