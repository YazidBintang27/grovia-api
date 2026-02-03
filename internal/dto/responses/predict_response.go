package responses

import (
	"grovia/internal/models"
	"time"
)

type PredictResponse struct {
	ID                int       `json:"id"`
	ToddlerID         int       `json:"toddlerID"`
	CreatedByID       int       `json:"createdByID"`
	Name              string    `json:"name"`
	Height            float64   `json:"height"`
	Age               int       `json:"age"`
	Sex               string    `json:"sex"`
	Zscore            float64   `json:"zscore"`
	NutritionalStatus string    `json:"nutritionalStatus"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func FromModelPredict(predict models.Predict) *PredictResponse {
	return &PredictResponse{
		ID:                predict.ID,
		ToddlerID:         predict.ToddlerID,
		CreatedByID:       predict.CreatedByID,
		Name:              predict.Name,
		Height:            predict.Height,
		Age:               predict.Age,
		Sex:               predict.Sex,
		Zscore:            predict.Zscore,
		NutritionalStatus: predict.NutritionalStatus,
		CreatedAt:         predict.CreatedAt,
		UpdatedAt:         predict.UpdatedAt,
	}
}

func FromModelPredictList(predicts []models.Predict) []PredictResponse {
	var predictResponse []PredictResponse
	for _, v := range predicts {
		predictResponse = append(predictResponse, *FromModelPredict(v))
	}
	return predictResponse
}
