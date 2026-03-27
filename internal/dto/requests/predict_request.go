package requests

import (
	"grovia/internal/models"
	"time"
)

type CreateIndividualPredictRequest struct {
	Height            float64 `json:"height" validate:"required,height"`
	Age               int     `json:"age" validate:"required,age"`
	Sex               string  `json:"sex" validate:"required"`
	NutritionalStatus string  `json:"nutritionalStatus" validate:"required"`
}

type UpdatePredictRequest struct {
	Height            *float64 `json:"height,omitempty" validate:"omitempty,height"`
	Age               *int     `json:"age,omitempty" validate:"omitempty,age"`
	Sex               *string  `json:"sex,omitempty" validate:"omitempty,oneof=male female"`
	Zscore            *float64 `json:"zscore,omitempty" validate:"omitempty"`
	NutritionalStatus *string  `json:"nutritionalStatus,omitempty" validate:"omitempty"`
}

// ToModel: Digunakan untuk mapping hasil dari ML API ke database model
func (r *CreateIndividualPredictRequest) ToModel(userID, locationID, age int, zscore float64) models.Predict {
	return models.Predict{
		CreatedByID:       userID,
		Name:              "", // Bisa diisi manual di service jika perlu
		Height:            r.Height,
		Age:               age,
		Sex:               r.Sex,
		Zscore:            zscore,
		NutritionalStatus: r.NutritionalStatus,
		LocationID:        locationID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

// ToUpdateModel: Mengamankan pointer dereference saat update
func (r *UpdatePredictRequest) ToUpdateModel(id int) models.Predict {
	predictModel := models.Predict{
		ID:        id,
		UpdatedAt: time.Now(),
	}

	if r.Height != nil {
		predictModel.Height = *r.Height
	}
	if r.Age != nil {
		predictModel.Age = *r.Age
	}
	if r.Sex != nil {
		predictModel.Sex = *r.Sex
	}
	if r.Zscore != nil {
		predictModel.Zscore = *r.Zscore
	}
	if r.NutritionalStatus != nil {
		predictModel.NutritionalStatus = *r.NutritionalStatus
	}

	return predictModel
}
