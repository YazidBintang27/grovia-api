package requests

type CreateIndividualPredictRequest struct {
	Height            float64 `json:"height" validate:"required"`
	Age               int     `json:"age" validate:"required"`
	Sex               string  `json:"sex" validate:"required"`
	NutritionalStatus string  `json:"nutritionalStatus" validate:"required"`
}

type UpdatePredictRequest struct {
	Height            *float64 `json:"height,omitempty"`
	Age               *int     `json:"age,omitempty"`
	Sex               string   `json:"sex,omitempty"`
	Zscore            *float64 `json:"zscore,omitempty"`
	NutritionalStatus *string  `json:"nutritionalStatus,omitempty"`
}
