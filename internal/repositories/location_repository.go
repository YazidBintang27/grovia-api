package repositories

import (
	"grovia/internal/models"
	"strings"

	"gorm.io/gorm"
)

type LocationRepository interface {
	CreateLocation(location *models.Location) (*models.Location, error)
	GetAllLocation(name string, limit, offset int) ([]models.Location, int, error)
	GetLocationByID(id int) (*models.Location, error)
	UpdateLocationByID(id int, location *models.Location) (*models.Location, error)
	DeleteLocationByID(id, userID int) error
}

type locationRepository struct {
	db *gorm.DB
}

// CreateLocation implements LocationRepository.
func (l *locationRepository) CreateLocation(location *models.Location) (*models.Location, error) {
	if err := l.db.Create(location).Error; err != nil {
		return nil, err
	}
	return location, nil
}

// DeleteLocationByID implements LocationRepository.
func (l *locationRepository) DeleteLocationByID(id, userID int) error {
	res := l.db.Model(&models.Location{}).
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

	return nil
}

// GetAllLocation implements LocationRepository.
func (l *locationRepository) GetAllLocation(name string, limit, offset int) ([]models.Location, int, error) {
	var locations []models.Location
	var total int64

	db := l.db.Model(&locations)

	if strings.TrimSpace(name) != "" {
		normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", ""))
		db = db.Where("REPLACE(LOWER(name), ' ', '') LIKE ?", "%"+normalizedName+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&locations).Error; err != nil {
		return nil, 0, err
	}

	return locations, int(total), nil
}

// GetLocationByID implements LocationRepository.
func (l *locationRepository) GetLocationByID(id int) (*models.Location, error) {
	var location models.Location

	tx := l.db.Where("id = ?", id).Find(&location)

	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &location, nil
}

// UpdateLocationByID implements LocationRepository.
func (l *locationRepository) UpdateLocationByID(id int, location *models.Location) (*models.Location, error) {
	if err := l.db.Model(&models.Location{}).Where("id = ?", id).Updates(location).Error; err != nil {
		return nil, err
	}

	var locationResponse models.Location

	tx := l.db.Where("id = ?", id).Find(&locationResponse)

	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &locationResponse, nil
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db: db}
}
