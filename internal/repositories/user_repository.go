package repositories

import (
	"grovia/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUser(id int) (*models.User, error)
	GetAllUser() ([]models.User, error)
	UpdateUser(id int, user *models.User) (*models.User, error)
	DeleteUser(id int) error
	FindRoleById(id int) (string, error)
	FindUsersByRole(roles []string) ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// FindUsersByRole implements UserRepository.
func (u *userRepository) FindUsersByRole(roles []string) ([]models.User, error) {
	var users []models.User
	err := u.db.Where("role IN ?", roles).Find(&users).Error
	return users, err
}

// GetAllUser implements UserRepository.
func (u *userRepository) GetAllUser() ([]models.User, error) {
	var users []models.User

	err := u.db.Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// FindRoleById implements UserRepository.
func (u *userRepository) FindRoleById(id int) (string, error) {
	var role string
	if err := u.db.Model(&models.User{}).
		Where("id = ?", id).
		Pluck("role", &role).Error; err != nil {
		return "", err
	}
	return role, nil
}

// CreateUser implements UserRepository.
func (u *userRepository) CreateUser(user *models.User) (*models.User, error) {
	err := u.db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser implements UserRepository.
func (u *userRepository) DeleteUser(id int) error {
	return u.db.Delete(&models.User{}, id).Error
}

// GetUser implements UserRepository.
func (u *userRepository) GetUser(id int) (*models.User, error) {
	var user models.User

	err := u.db.First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser implements UserRepository.
func (u *userRepository) UpdateUser(id int, user *models.User) (*models.User, error) {
	var existing models.User

	if err := u.db.First(&existing, id).Error; err != nil {
		return nil, err
	}

	err := u.db.Model(&existing).Updates(user).Error
	if err != nil {
		return nil, err
	}

	return &existing, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
