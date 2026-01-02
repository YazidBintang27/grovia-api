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

// DeleteToddlerByID implements ToddlerRepository.
func (t *toddlerRepository) DeleteToddlerByID(id int, locationID, userID int) error {
	db := t.db.Model(&models.Toddler{})

	if locationID == 1 {
		res := db.Where("id = ?", id).Updates(map[string]any{
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
		res := db.Where("id = ? AND location_id = ?", id, locationID).Updates(map[string]any{
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
	var toddler models.Toddler

	db := t.db.Where("id = ? AND deleted_at IS NULL", id)

	if locationID != 1 {
		db = db.Where("location_id = ?", locationID)
	}

	if err := db.First(&toddler).Error; err != nil {
		return nil, err
	}

	return &toddler, nil
}

// UpdateToddlerByID implements ToddlerRepository.
func (t *toddlerRepository) UpdateToddlerByID(id int, locationID int, toddler *models.Toddler) (*models.Toddler, error) {
	db := t.db.Model(toddler).Where("id = ? AND deleted_at IS NULL", id)

	if locationID != 1 {
		db = db.Where("location_id = ?", locationID)
	}

	if err := db.Updates(toddler).Error; err != nil {
		return nil, err
	}

	if db.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var toddlerResponse models.Toddler
	dbFetch := t.db.Where("id = ?", id)

	if locationID != 1 {
		dbFetch = dbFetch.Where("location_id = ?", locationID)
	}

	if err := dbFetch.First(&toddlerResponse).Error; err != nil {
		return nil, err
	}

	return &toddlerResponse, nil
}

func NewToddlerRepository(db *gorm.DB) ToddlerRepository {
	return &toddlerRepository{db: db}
}
