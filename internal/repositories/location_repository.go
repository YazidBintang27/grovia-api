package repositories

import (
	"grovia/internal/models"
	"strings"

	"gorm.io/gorm"
)

type LocationRepository interface {
	CreateLocation(location *models.Location) (*models.Location, error)
	GetAllLocation(name string) ([]models.Location, error)
	GetLocationByID(id int) (*models.Location, error)
	UpdateLocationByID(id int, location *models.Location) (*models.Location, error)
	DeleteLocationByID(id int) error
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
func (l *locationRepository) DeleteLocationByID(id int) error {
	tx := l.db.Where("id = ?", id).Delete(&models.Location{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetAllLocation implements LocationRepository.
func (l *locationRepository) GetAllLocation(name string) ([]models.Location, error) {
	var locations []models.Location

	db := l.db.Model(&locations)

	if strings.TrimSpace(name) != "" {
		normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", ""))
		db = db.Where("REPLACE(LOWER(name), ' ', '') LIKE ?", "%"+normalizedName+"%")
	}

	err := db.Find(&locations).Error
	return locations, err
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
