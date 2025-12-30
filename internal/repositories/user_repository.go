package repositories

import (
	"grovia/internal/models"
	"strings"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUser(id int) (*models.User, error)
	GetAllUser() ([]models.User, error)
	UpdateUser(id int, user *models.User) (*models.User, error)
	DeleteUser(id int) error
	FindRoleById(id int) (string, error)
	FindUsersByRole(roles []string, name string, locationID, limit, offset int) ([]models.User, int, error)
}

type userRepository struct {
	db *gorm.DB
}

// FindUsersByRole implements UserRepository.
func (u *userRepository) FindUsersByRole(roles []string, name string, locationID, limit, offset int) ([]models.User, int, error) {
	var users []models.User
	var total int64
	db := u.db.Model(&models.User{})

	db = db.Where("is_active = true AND role IN ?", roles)

	if locationID != 1 {
		db = db.Where("location_id = ?", locationID)
	}

	if strings.TrimSpace(name) != "" {
		normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", ""))
		db = db.Where("REPLACE(LOWER(name), ' ', '') LIKE ?", "%"+normalizedName+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

// GetAllUser implements UserRepository.
func (u *userRepository) GetAllUser() ([]models.User, error) {
	var users []models.User

	err := u.db.Where("is_active = true").Order("created_at DESC").Find(&users).Error

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
	tx := u.db.Model(&models.User{}).Where("id = ? AND is_active = true", id).Update("is_active", false)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetUser implements UserRepository.
func (u *userRepository) GetUser(id int) (*models.User, error) {
	var user models.User

	tx := u.db.First(&user, id)

	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &user, nil
}

// UpdateUser implements UserRepository.
func (u *userRepository) UpdateUser(id int, user *models.User) (*models.User, error) {
	var existing models.User

	if err := u.db.First(&existing, id).Error; err != nil {
		return nil, err
	}

	tx := u.db.Model(&existing).Updates(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &existing, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
