package repositories

import (
	"errors"
	"grovia/internal/models"
	"strings"

	"gorm.io/gorm"
)

type ParentRepository interface {
	CreateParent(parent *models.Parent) (*models.Parent, error)
	GetAllParent(locationID, limit, offset int, name string) ([]models.Parent, int, error)
	GetParentByID(id, locationID int) (*models.Parent, error)
	UpdateParentByID(id, locationID int, parent *models.Parent) (*models.Parent, error)
	DeleteParentByID(id, locationID, userID int) error
	FindParentByPhoneNumber(phoneNumber string) (*models.Parent, error)
	GetAllParentAllLocation(name string, limit, offset int) ([]models.Parent, int, error)
}

type parentRepository struct {
	db *gorm.DB
}

// GetAllParentAllLocation implements ParentRepository.
func (p *parentRepository) GetAllParentAllLocation(name string, limit, offset int) ([]models.Parent, int, error) {
	var parents []models.Parent
	var total int64

	db := p.db.Model(&parents).Where("deleted_at IS NULL")

	if strings.TrimSpace(name) != "" {
		normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", ""))
		db = db.Where("REPLACE(LOWER(name), ' ', '') LIKE ?", "%"+normalizedName+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&parents).Error; err != nil {
		return nil, 0, err
	}
	return parents, int(total), nil
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
func (p *parentRepository) DeleteParentByID(id, locationID, userID int) error {
	if locationID == 1 {
		res := p.db.Model(&models.Parent{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"deleted_by_id": userID,
				"deleted_at":    gorm.Expr("NOW()"),
			})

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
	} else {
		res := p.db.Model(&models.Parent{}).
			Where("id = ? AND location_id = ?", id, locationID).
			Updates(map[string]any{
				"deleted_by_id": userID,
				"deleted_at":    gorm.Expr("NOW()"),
			})

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
	}

	if err := p.db.Model(&models.Toddler{}).
		Where("parent_id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_by_id": userID,
			"deleted_at":    gorm.Expr("NOW()"),
		}).Error; err != nil {
		return err
	}
	
	return nil
}

// GetAllParent implements ParentRepository.
func (p *parentRepository) GetAllParent(locationID, limit, offset int, name string) ([]models.Parent, int, error) {
	var parents []models.Parent
	var total int64

	db := p.db.Model(&parents).Where("deleted_at IS NULL")

	db = db.Where("location_id = ?", locationID)

	if strings.TrimSpace(name) != "" {
		normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", ""))
		db = db.Where("REPLACE(LOWER(name), ' ', '') LIKE ?", "%"+normalizedName+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&parents).Error; err != nil {
		return nil, 0, err
	}

	return parents, int(total), nil
}

// GetParentByID implements ParentRepository.
func (p *parentRepository) GetParentByID(id, locationID int) (*models.Parent, error) {
	var parent models.Parent

	db := p.db.Preload("Toddlers", func(tx *gorm.DB) *gorm.DB {
		return tx.Where("deleted_at IS NULL")
	})

	if locationID == 1 {
		tx := db.Where("id = ?", id).Find(&parent)

		if tx.Error != nil {
			return nil, tx.Error
		}
		if tx.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	} else {
		tx := db.Where("id = ? AND location_id = ?", id, locationID).Find(&parent)

		if tx.Error != nil {
			return nil, tx.Error
		}
		if tx.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	}
	return &parent, nil
}

// UpdateParentByID implements ParentRepository.
func (p *parentRepository) UpdateParentByID(id, locationID int, parent *models.Parent) (*models.Parent, error) {

	if locationID == 1 {
		err := p.db.Model(parent).Where("id = ?", id).Updates(parent).Error

		if err != nil {
			return nil, err
		}
	} else {
		err := p.db.Model(parent).Where("id = ? AND location_id = ?", id, locationID).Updates(parent).Error

		if err != nil {
			return nil, err
		}
	}

	var parentResponse models.Parent

	if locationID == 1 {
		tx := p.db.Where("id = ?", id).First(&parentResponse)

		if tx.Error != nil {
			return nil, tx.Error
		}
		if tx.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	} else {
		tx := p.db.Where("id = ? AND location_id = ?", id, locationID).First(&parentResponse)

		if tx.Error != nil {
			return nil, tx.Error
		}
		if tx.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	}
	return &parentResponse, nil
}

func NewParentRepository(db *gorm.DB) ParentRepository {
	return &parentRepository{db: db}
}
