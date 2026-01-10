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
	LocationID        int       `json:"locationID" validate:"required"`
	PhoneNumber       string    `json:"phoneNumber" validate:"required,phone"`
}

type UpdateToddlerRequest struct {
	Name              *string               `form:"name,omitempty" validate:"omitempty"`
	Birthdate         *time.Time            `form:"birthdate,omitempty" validate:"omitempty"`
	Sex               string                `form:"sex,omitempty" validate:"omitempty,oneof=male female"`
	Height            *float64              `form:"height,omitempty" validate:"omitempty,height"`
	ProfilePicture    *multipart.FileHeader `form:"profilePicture,omitempty"`
	NutritionalStatus *string               `form:"nutritionalStatus,omitempty" validate:"omitempty"`
	LocationID        *int                  `form:"locationID,omitempty" validate:"omitempty"`
	PhoneNumber       *string               `form:"phoneNumber,omitempty" validate:"omitempty,phone"`
}

type CreateToddlerWithParentRequest struct {
	Toddler CreateToddlerRequest `json:"toddler"`
	Parent  CreateParentRequest  `json:"parent"`
}
