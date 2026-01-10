package requests

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
