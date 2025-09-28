package repositories

import (
	"grovia/internal/models"

	"gorm.io/gorm"
)

type ParentRepository interface {
	CreateParentWithToddlers(parent *models.Parent, toddlers []models.Toddler) (*models.Parent, error)
	GetAllParent(locationID int) ([]models.Parent, error)
	GetParentByID(id, locationID int) (*models.Parent, error)
	UpdateParentByID(id, locationID int, parent *models.Parent) (*models.Parent, error)
	DeleteParentByID(id, locationID int) error
}

type parentRepository struct {
	db *gorm.DB
}

// CreateParentWithToddlers implements ParentRepository.
func (p *parentRepository) CreateParentWithToddlers(parent *models.Parent, toddlers []models.Toddler) (*models.Parent, error) {
	// create parent dulu
	if err := p.db.Create(parent).Error; err != nil {
		return nil, err
	}

	for i := range toddlers {
		toddlers[i].ParentID = parent.ID
		if err := p.db.Create(&toddlers[i]).Error; err != nil {
			return nil, err
		}
	}

	return parent, nil
}

// DeleteParentByID implements ParentRepository.
func (p *parentRepository) DeleteParentByID(id, locationID int) error {
	return p.db.Where("id = ? AND location_id = ?", id, locationID).Delete(&models.Parent{}).Error
}

// GetAllParent implements ParentRepository.
func (p *parentRepository) GetAllParent(locationID int) ([]models.Parent, error) {
	var parents []models.Parent

	err := p.db.Where("location_id = ?", locationID).First(&parents).Error

	if err != nil {
		return nil, err
	}

	return parents, nil
}

// GetParentByID implements ParentRepository.
func (p *parentRepository) GetParentByID(id, locationID int) (*models.Parent, error) {
	var parent models.Parent

	err := p.db.Where("id = ? AND location_id = ?", id, locationID).Find(&parent).Error

	if err != nil {
		return nil, err
	}

	return &parent, nil
}

// UpdateParentByID implements ParentRepository.
func (p *parentRepository) UpdateParentByID(id, locationID int, parent *models.Parent) (*models.Parent, error) {
	err := p.db.Model(&models.Parent{}).Where("id = ? AND location_id = ?", id, locationID).Updates(parent).Error

	if err != nil {
		return nil, err
	}

	var parentResponse models.Parent
	err = p.db.Where("id = ? AND location_id = ?", id, locationID).First(&parentResponse).Error

	if err != nil {
		return nil, err
	}

	return &parentResponse, nil
}

func NewParentRepository(db *gorm.DB) ParentRepository {
	return &parentRepository{db: db}
}
