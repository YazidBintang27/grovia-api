package repositories

import (
	"grovia/internal/models"

	"gorm.io/gorm"
)

type PredictRepository interface {
	CreateIndividualPredict(predict *models.Predict, locationID, toddlerID int) (*models.Predict, error)
	GetAllPredict(locationID, limit, offset int) ([]models.Predict, int, error)
	GetAllPredictByToddlerID(locationID, toddlerID int) ([]models.Predict, error)
	GetPredictByID(id int) (*models.Predict, error)
	UpdatePredictByID(id int, predict *models.Predict) (*models.Predict, error)
	DeletePredictByID(id, locationID, userID int) error
	GetAllPredictAllLocation(limit, offset int) ([]models.Predict, int, error)
}

type predictRepository struct {
	db *gorm.DB
}

// GetAllPredictAllLocation implements PredictRepository.
func (p *predictRepository) GetAllPredictAllLocation(limit, offset int) ([]models.Predict, int, error) {
	var predicts []models.Predict
	var total int64

	db := p.db.Model(&predicts).Where("deleted_at IS NULL")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&predicts).Error; err != nil {
		return nil, 0, err
	}

	return predicts, int(total), nil
}

// CreateIndividualPredict implements PredictRepository.
func (p *predictRepository) CreateIndividualPredict(predict *models.Predict, locationID, toddlerID int) (*models.Predict, error) {
	predict.ToddlerID = toddlerID
	if err := p.db.Create(predict).Error; err != nil {
		return nil, err
	}

	return predict, nil
}

// DeletePredictByID implements PredictRepository.
func (p *predictRepository) DeletePredictByID(id int, locationID, userID int) error {
	db := p.db.Model(&models.Predict{}).Where("id = ? AND deleted_at IS NULL", id)

	if locationID != 1 {
		db = db.Where("location_id = ?", locationID)
	}

	res := db.Updates(map[string]any{
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

// GetAllPredict implements PredictRepository.
func (p *predictRepository) GetAllPredict(locationID, limit, offset int) ([]models.Predict, int, error) {
	var predicts []models.Predict
	var total int64

	db := p.db.Model(&predicts).Where("location_id = ? AND deleted_at IS NULL", locationID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&predicts).Error; err != nil {
		return nil, 0, err
	}

	return predicts, int(total), nil
}

// GetAllPredictByToddlerID implements PredictRepository.
func (p *predictRepository) GetAllPredictByToddlerID(locationID, toddlerID int) ([]models.Predict, error) {
	var predicts []models.Predict

	db := p.db.Where("toddler_id = ? AND deleted_at IS NULL", toddlerID)

	if locationID != 1 {
		db = db.Where("location_id = ?", locationID)
	}

	if err := db.Order("created_at DESC").Find(&predicts).Error; err != nil {
		return nil, err
	}

	return predicts, nil
}

// GetPredictByID implements PredictRepository.
func (p *predictRepository) GetPredictByID(id int) (*models.Predict, error) {
	var predict models.Predict

	if err := p.db.Where("id = ? AND deleted_at IS NULL", id).First(&predict).Error; err != nil {
		return nil, err
	}

	return &predict, nil
}

// UpdatePredictByID implements PredictRepository.
func (p *predictRepository) UpdatePredictByID(id int, predict *models.Predict) (*models.Predict, error) {
	res := p.db.Model(&models.Predict{}).Where("id = ? AND deleted_at IS NULL", id).Updates(predict)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var predictResponse models.Predict
	if err := p.db.Where("id = ?", id).First(&predictResponse).Error; err != nil {
		return nil, err
	}

	return &predictResponse, nil
}

func NewPredictRepository(db *gorm.DB) PredictRepository {
	return &predictRepository{db: db}
}