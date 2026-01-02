package responses

import "time"

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
