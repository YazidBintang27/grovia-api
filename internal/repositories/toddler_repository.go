package repositories

import (
	"grovia/internal/models"

	"gorm.io/gorm"
)

type ToddlerRepository interface {
	GetAllToddler(locationID int) ([]models.Toddler, error)
	GetToddlerByID(id, locationID int) (*models.Toddler, error)
	UpdateToddlerByID(id, locationID int, toddler *models.Toddler) (*models.Toddler, error)
	DeleteToddlerByID(id, locationID int) error
	FindParentIDByID(id, locationID int) (*int, error)
}

type toddlerRepository struct {
	db *gorm.DB
}

// FindParentIDByID implements ToddlerRepository.
func (t *toddlerRepository) FindParentIDByID(id int, locationID int) (*int, error) {
	var parentID *int
	if err := t.db.Model(&models.Toddler{}).
		Where("id = ? AND location_id = ?", id, locationID).
		Pluck("parent_id", &parentID).Error; err != nil {
		return nil, err
	}
	return parentID, nil
}

// DeleteToddlerByID implements ToddlerRepository.
func (t *toddlerRepository) DeleteToddlerByID(id int, locationID int) error {
	parentID, err := t.FindParentIDByID(id, locationID)
	if err != nil {
		return err
	}
	return t.db.Where("id = ? AND location_id = ? AND parent_id = ?", id, locationID, parentID).Delete(&models.Toddler{}).Error
}

// GetAllToddler implements ToddlerRepository.
func (t *toddlerRepository) GetAllToddler(locationID int) ([]models.Toddler, error) {
	var toddlers []models.Toddler

	err := t.db.Where("location_id = ?", locationID).First(&toddlers).Error

	if err != nil {
		return nil, err
	}

	return toddlers, nil
}

// GetToddlerByID implements ToddlerRepository.
func (t *toddlerRepository) GetToddlerByID(id int, locationID int) (*models.Toddler, error) {

	parentID, err := t.FindParentIDByID(id, locationID)
	if err != nil {
		return nil, err
	}

	var toddler models.Toddler

	err = t.db.Where("id = ? AND location_id = ? AND parent_id = ?", id, locationID, parentID).Find(&toddler).Error

	if err != nil {
		return nil, err
	}

	return &toddler, nil
}

// UpdateToddlerByID implements ToddlerRepository.
func (t *toddlerRepository) UpdateToddlerByID(id int, locationID int, toddler *models.Toddler) (*models.Toddler, error) {

	parentID, err := t.FindParentIDByID(id, locationID)
	if err != nil {
		return nil, err
	}

	err = t.db.Model(&toddler).Where("id = ? AND location_id = ? AND parent_id = ?", id, locationID, parentID).Updates(toddler).Error

	if err != nil {
		return nil, err
	}

	var toddlerResponse models.Toddler

	err = t.db.Where("id = ? AND location_id = ?", id, locationID).First(&toddlerResponse).Error

	if err != nil {
		return nil, err
	}

	return &toddlerResponse, nil
}

func NewToddlerRepository(db *gorm.DB) ToddlerRepository {
	return &toddlerRepository{db: db}
}
