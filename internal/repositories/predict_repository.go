package repositories

import (
	"grovia/internal/models"

	"gorm.io/gorm"
)

type PredictRepository interface {
	CreateIndividualPredict(predict *models.Predict, locationID, toddlerID int) (*models.Predict, error)
	GetAllPredict(locationID int) ([]models.Predict, error)
	GetAllPredictByToddlerID(locationID, toddlerID int) ([]models.Predict, error)
	GetPredictByID(id int) (*models.Predict, error)
	UpdatePredictByID(id int, predict *models.Predict) (*models.Predict, error)
	DeletePredictByID(id, locationID int) error
	FindToddlerIDByID(id, locationID int) (*int, error)
	GetAllPredictAllLocation() ([]models.Predict, error)
}

type predictRepository struct {
	db *gorm.DB
}

// GetAllPredictAllLocation implements PredictRepository.
func (p *predictRepository) GetAllPredictAllLocation() ([]models.Predict, error) {
	var predicts []models.Predict

	err := p.db.Find(&predicts).Error

	if err != nil {
		return nil, err
	}

	return predicts, nil
}

// FindToddlerIDByID implements PredictRepository.
func (p *predictRepository) FindToddlerIDByID(id int, locationID int) (*int, error) {
	var toddler models.Toddler
	err := p.db.Select("id").Where("id = ? AND location_id = ?", id, locationID).First(&toddler).Error
	if err != nil {
		return nil, err
	}
	return &toddler.ID, nil
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
func (p *predictRepository) DeletePredictByID(id int, locationID int) error {
	toddlerID, err := p.FindToddlerIDByID(id, locationID)

	if err != nil {
		return err
	}

	tx := p.db.Where("id = ? AND location_id = ? AND toddler_id = ?", id, locationID, toddlerID)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetAllPredict implements PredictRepository.
func (p *predictRepository) GetAllPredict(locationID int) ([]models.Predict, error) {
	var predicts []models.Predict

	err := p.db.Where("location_id = ?", locationID).Find(&predicts).Error

	if err != nil {
		return nil, err
	}

	return predicts, nil
}

// GetAllPredictByToddlerID implements PredictRepository.
func (p *predictRepository) GetAllPredictByToddlerID(locationID, toddlerID int) ([]models.Predict, error) {
	var predicts []models.Predict

	tx := p.db.Where("location_id = ? AND toddler_id = ?", locationID, toddlerID).First(&predicts)

	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return predicts, nil
}

// GetPredictByID implements PredictRepository.
func (p *predictRepository) GetPredictByID(id int) (*models.Predict, error) {
	var predict models.Predict

	tx := p.db.Where("id = ?", id).Find(&predict)

	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &predict, nil
}

// UpdatePredictByID implements PredictRepository.
func (p *predictRepository) UpdatePredictByID(id int, predict *models.Predict) (*models.Predict, error) {

	tx := p.db.Model(&predict).Where("id = ?", id).Updates(predict)

	var predictResponse models.Predict

	err := p.db.Where("id = ?", id).First(&predictResponse).Error

	if err != nil {
		return nil, err
	}

	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &predictResponse, nil
}

func NewPredictRepository(db *gorm.DB) PredictRepository {
	return &predictRepository{db: db}
}
