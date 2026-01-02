package repositories

import (
	"grovia/internal/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByPhoneNumber(phoneNumber string) (*models.User, error)
	ResetPassword(phoneNumber, newPassword string) error
	FindByID(id int) (*models.User, error)
}

type authRepository struct {
	db *gorm.DB
}

// FindByID implements AuthRepository.
func (a *authRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := a.db.Where("is_active = true").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ResetPassword implements AuthRepository.
func (a *authRepository) ResetPassword(phoneNumber string, newPassword string) error {
	return a.db.Model(&models.User{}).
		Where("phone_number = ?", phoneNumber).
		Update("password", newPassword).
		Error
}

// FindByPhoneNumber implements AuthRepository.
func (a *authRepository) FindByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	err := a.db.Where("phone_number = ? AND is_active = true", phoneNumber).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}
