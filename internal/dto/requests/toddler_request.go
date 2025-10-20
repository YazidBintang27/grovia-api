package requests

import (
	"mime/multipart"
	"time"
)

type CreateToddlerRequest struct {
	Name              string    `json:"name" validate:"required"`
	Birthdate         time.Time `json:"birthdate" validate:"required"`
	Sex               string    `json:"sex" validate:"required,oneof=male female"`
	Height            float64   `json:"height" validate:"required"`
	NutritionalStatus string    `json:"nutritionalStatus" validate:"required"`
	LocationID        int       `json:"locationID"`
	PhoneNumber       string    `json:"phoneNumber" validate:"required"`
}

type UpdateToddlerRequest struct {
	Name              *string               `form:"name,omitempty"`
	Birthdate         *time.Time            `form:"birthdate,omitempty"`
	Sex               string                `form:"sex,omitempty"`
	Height            *float64              `form:"height,omitempty"`
	ProfilePicture    *multipart.FileHeader `form:"profilePicture,omitempty"`
	NutritionalStatus *string               `form:"nutritionalStatus,omitempty"`
	LocationID        *int                  `form:"locationID,omitempty"`
	PhoneNumber       *string               `form:"phoneNumber,omitempty"`
}

type CreateToddlerWithParentRequest struct {
	Toddler CreateToddlerRequest `json:"toddler"`
	Parent  CreateParentRequest  `json:"parent"`
}
