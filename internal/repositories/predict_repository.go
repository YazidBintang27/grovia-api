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
	FindToddlerIDByID(id, locationID int) (*int, error)
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

	tx := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&predicts)

	if tx.Error != nil {
		return nil, 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	}

	return predicts, int(total), nil
}

// FindToddlerIDByID implements PredictRepository.
func (p *predictRepository) FindToddlerIDByID(id int, locationID int) (*int, error) {
	var predict models.Predict
	if locationID == 1 {
		err := p.db.Select("toddler_id").Where("id = ? ", id).First(&predict).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := p.db.Select("toddler_id").Where("id = ? AND location_id = ?", id, locationID).First(&predict).Error
		if err != nil {
			return nil, err
		}
	}

	return &predict.ToddlerID, nil
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
	toddlerID, err := p.FindToddlerIDByID(id, locationID)

	if err != nil {
		return err
	}

	db := p.db.Model(&models.Predict{})

	if locationID == 1 {
		res := db.Where("id = ? AND toddler_id = ?", id, toddlerID).Updates(map[string]any{
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
		res := db.Where("id = ? AND toddler_id = ? AND location_id = ?", id, toddlerID, locationID).Updates(map[string]any{
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

	return nil
}

// GetAllPredict implements PredictRepository.
func (p *predictRepository) GetAllPredict(locationID, limit, offset int) ([]models.Predict, int, error) {
	var predicts []models.Predict
	var total int64
	db := p.db.Model(&predicts).Where("deleted_at IS NULL")

	tx := db.Where("location_id = ?", locationID).Limit(limit).Offset(offset).Order("created_at DESC").Find(&predicts)

	if tx.Error != nil {
		return nil, 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	}

	return predicts, int(total), nil
}

// GetAllPredictByToddlerID implements PredictRepository.
func (p *predictRepository) GetAllPredictByToddlerID(locationID, toddlerID int) ([]models.Predict, error) {
	var predicts []models.Predict

	if locationID == 1 {
		tx := p.db.Where("toddler_id = ? AND deleted_at IS NULL", toddlerID).Order("created_at DESC").Find(&predicts)

		if tx.Error != nil {
			return nil, tx.Error
		}
		if tx.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	} else {
		tx := p.db.Where("location_id = ? AND toddler_id = ? AND deleted_at IS NULL", locationID, toddlerID).Order("created_at DESC").Find(&predicts)

		if tx.Error != nil {
			return nil, tx.Error
		}
		if tx.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
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
