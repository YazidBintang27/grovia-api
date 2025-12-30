package repositories

import (
	"errors"
	"grovia/internal/models"
	"strings"

	"gorm.io/gorm"
)

type ToddlerRepository interface {
	CreateToddler(toddler *models.Toddler) (*models.Toddler, error)
	GetAllToddler(locationID, limit, offset int, name string) ([]models.Toddler, int, error)
	GetToddlerByID(id, locationID int) (*models.Toddler, error)
	UpdateToddlerByID(id, locationID int, toddler *models.Toddler) (*models.Toddler, error)
	DeleteToddlerByID(id, locationID, userID int) error
	FindParentIDByID(id, locationID int) (*int, error)
	FindToddlerByName(parentID int, name string) (bool, *models.Toddler, error)
	GetAllToddlerAllLocation(name string, limit, offset int) ([]models.Toddler, int, error)
}

type toddlerRepository struct {
	db *gorm.DB
}

// GetAllToddlerAllLocation implements ToddlerRepository.
func (t *toddlerRepository) GetAllToddlerAllLocation(
	name string,
	limit, offset int,
) ([]models.Toddler, int, error) {

	var toddlers []models.Toddler
	var total int64

	db := t.db.
		Model(&models.Toddler{}).
		Joins("LEFT JOIN parents ON parents.id = toddlers.parent_id").
		Where("toddlers.deleted_at IS NULL")

	if strings.TrimSpace(name) != "" {
		normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", ""))

		db = db.Where(`
			REPLACE(LOWER(toddlers.name), ' ', '') LIKE ?
			OR
			REPLACE(LOWER(parents.name), ' ', '') LIKE ?
		`, "%"+normalizedName+"%", "%"+normalizedName+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.
		Preload("Parent").
		Limit(limit).
		Offset(offset).
		Order("toddlers.created_at DESC").
		Find(&toddlers).Error; err != nil {
		return nil, 0, err
	}

	return toddlers, int(total), nil
}

// FindToddlerByName implements ToddlerRepository.
func (t *toddlerRepository) FindToddlerByName(parentID int, name string) (bool, *models.Toddler, error) {
	normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", ""))

	var toddler models.Toddler
	err := t.db.Where("parent_id = ? AND REPLACE(LOWER(name), ' ', '') = ?", parentID, normalizedName).
		First(&toddler).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, err
		}
		return false, nil, err
	}
	return true, &toddler, nil
}

// CreateToddler implements ToddlerRepository.
func (t *toddlerRepository) CreateToddler(toddler *models.Toddler) (*models.Toddler, error) {
	if err := t.db.Create(toddler).Error; err != nil {
		return nil, err
	}

	return toddler, nil
}

// FindParentIDByID implements ToddlerRepository.
func (t *toddlerRepository) FindParentIDByID(id int, locationID int) (*int, error) {
	var parentID int
	if locationID == 1 {
		if err := t.db.Model(&models.Toddler{}).
			Select("parent_id").
			Where("id = ?", id).
			Take(&parentID).Error; err != nil {
			return nil, err
		}
	} else {
		if err := t.db.Model(&models.Toddler{}).
			Select("parent_id").
			Where("id = ? AND location_id = ?", id, locationID).
			Take(&parentID).Error; err != nil {
			return nil, err
		}
	}
	return &parentID, nil
}

// DeleteToddlerByID implements ToddlerRepository.
func (t *toddlerRepository) DeleteToddlerByID(id int, locationID, userID int) error {
	parentID, err := t.FindParentIDByID(id, locationID)
	if err != nil {
		return err
	}

	db := t.db.Model(&models.Toddler{})

	if locationID == 1 {
		res := db.Where("id = ? AND parent_id = ?", id, parentID).Updates(map[string]any{
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
		res := db.Where("id = ? AND location_id = ? AND parent_id = ?", id, locationID, parentID).Updates(map[string]any{
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

	if err := t.db.Model(&models.Predict{}).
		Where("toddler_id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_by_id": userID,
			"deleted_at":    gorm.Expr("NOW()"),
		}).Error; err != nil {
		return err
	}

	return nil
}

// GetAllToddler implements ToddlerRepository.
func (t *toddlerRepository) GetAllToddler(locationID, limit, offset int, name string) ([]models.Toddler, int, error) {
	var toddlers []models.Toddler
	var total int64
	db := t.db.Model(&toddlers).Where("location_id = ? AND deleted_at IS NULL", locationID)

	if strings.TrimSpace(name) != "" {
		normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", ""))
		db = db.Where("REPLACE(LOWER(name), ' ', '') LIKE ?", "%"+normalizedName+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&toddlers).Error; err != nil {
		return nil, 0, err
	}

	return toddlers, int(total), nil
}

// GetToddlerByID implements ToddlerRepository.
func (t *toddlerRepository) GetToddlerByID(id int, locationID int) (*models.Toddler, error) {

	parentID, err := t.FindParentIDByID(id, locationID)
	if err != nil {
		return nil, err
	}

	var toddler models.Toddler

	if locationID == 1 {
		tx := t.db.Where("id = ? AND parent_id = ?", id, parentID).Find(&toddler)

		if tx.Error != nil {
			return nil, tx.Error
		}
		if tx.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	} else {
		tx := t.db.Where("id = ? AND location_id = ? AND parent_id = ?", id, locationID, parentID).Find(&toddler)

		if tx.Error != nil {
			return nil, tx.Error
		}
		if tx.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	}

	return &toddler, nil
}

// UpdateToddlerByID implements ToddlerRepository.
func (t *toddlerRepository) UpdateToddlerByID(id int, locationID int, toddler *models.Toddler) (*models.Toddler, error) {

	parentID, err := t.FindParentIDByID(id, locationID)
	if err != nil {
		return nil, err
	}

	if locationID == 1 {
		tx := t.db.Model(toddler).Where("id = ? AND parent_id = ?", id, parentID).Updates(toddler)

		if tx.Error != nil {
			return nil, tx.Error
		}
		if tx.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	} else {
		tx := t.db.Model(toddler).Where("id = ? AND location_id = ? AND parent_id = ?", id, locationID, parentID).Updates(toddler)

		if tx.Error != nil {
			return nil, tx.Error
		}
		if tx.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	}

	var toddlerResponse models.Toddler

	if locationID == 1 {
		tx := t.db.Where("id = ?", id).First(&toddlerResponse)

		if tx.Error != nil {
			return nil, tx.Error
		}
		if tx.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	} else {
		tx := t.db.Where("id = ? AND location_id = ?", id, locationID).First(&toddlerResponse)

		if tx.Error != nil {
			return nil, tx.Error
		}
		if tx.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	}

	return &toddlerResponse, nil
}

func NewToddlerRepository(db *gorm.DB) ToddlerRepository {
	return &toddlerRepository{db: db}
}
