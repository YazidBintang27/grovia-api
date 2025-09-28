package requests

import (
	"mime/multipart"
	"time"
)

type CreateToddlerRequest struct {
	Name      string    `json:"name" validate:"required"`
	Birthdate time.Time `json:"birthdate" validate:"required"`
	Gender    string    `json:"gender" validate:"required,oneof=male female"`
	Height    float64   `json:"height" validate:"required"`
}

type UpdateToddlerRequest struct {
	Name           string                `form:"name" validate:"required"`
	Birthdate      time.Time             `form:"birthdate" validate:"required"`
	Gender         string                `form:"gender" validate:"required,oneof=male female"`
	Height         float64               `form:"height" validate:"required"`
	ProfilePicture *multipart.FileHeader `form:"profilePicture"`
}
