package repositories

import (
	"errors"
	"grovia/internal/models"

	"gorm.io/gorm"
)

type ParentRepository interface {
	CreateParent(parent *models.Parent) (*models.Parent, error)
	GetAllParent(locationID int) ([]models.Parent, error)
	GetParentByID(id, locationID int) (*models.Parent, error)
	UpdateParentByID(id, locationID int, parent *models.Parent) (*models.Parent, error)
	DeleteParentByID(id, locationID int) error
	FindParentByPhoneNumber(phoneNumber string) (*models.Parent, error)
	GetAllParentAllLocation() ([]models.Parent, error)
}

type parentRepository struct {
	db *gorm.DB
}

// GetAllParentAllLocation implements ParentRepository.
func (p *parentRepository) GetAllParentAllLocation() ([]models.Parent, error) {
	var parents []models.Parent

	err := p.db.Find(&parents).Error

	if err != nil {
		return nil, err
	}

	return parents, nil
}

// CreateParent implements ParentRepository.
func (p *parentRepository) CreateParent(parent *models.Parent) (*models.Parent, error) {
	if err := p.db.Create(parent).Error; err != nil {
		return nil, err
	}
	return parent, nil
}

// FindParentByPhoneNumber implements ParentRepository.
func (p *parentRepository) FindParentByPhoneNumber(phoneNumber string) (*models.Parent, error) {
	var parent models.Parent
	err := p.db.Where("phone_number = ?", phoneNumber).First(&parent).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &parent, nil
}

// DeleteParentByID implements ParentRepository.
func (p *parentRepository) DeleteParentByID(id, locationID int) error {
	tx := p.db.Where("id = ? AND location_id = ?", id, locationID).Delete(&models.Parent{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetAllParent implements ParentRepository.
func (p *parentRepository) GetAllParent(locationID int) ([]models.Parent, error) {
	var parents []models.Parent

	err := p.db.Where("location_id = ?", locationID).Find(&parents).Error

	if err != nil {
		return nil, err
	}

	return parents, nil
}

// GetParentByID implements ParentRepository.
func (p *parentRepository) GetParentByID(id, locationID int) (*models.Parent, error) {
	var parent models.Parent

	tx := p.db.Where("id = ? AND location_id = ?", id, locationID).Find(&parent)

	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &parent, nil
}

// UpdateParentByID implements ParentRepository.
func (p *parentRepository) UpdateParentByID(id, locationID int, parent *models.Parent) (*models.Parent, error) {
	err := p.db.Model(parent).Where("id = ? AND location_id = ?", id, locationID).Updates(parent).Error

	if err != nil {
		return nil, err
	}

	var parentResponse models.Parent
	tx := p.db.Where("id = ? AND location_id = ?", id, locationID).First(&parentResponse)

	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &parentResponse, nil
}

func NewParentRepository(db *gorm.DB) ParentRepository {
	return &parentRepository{db: db}
}
