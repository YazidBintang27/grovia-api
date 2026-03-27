package requests

import (
	"grovia/internal/models"
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

// ToModel mengonversi CreateToddlerRequest menjadi models.Toddler
func (r *CreateToddlerRequest) ToModel(userID int, parentID int) models.Toddler {
	return models.Toddler{
		ParentID:          parentID,
		CreatedByID:       userID,
		UpdatedByID:       userID,
		Name:              r.Name,
		Birthdate:         r.Birthdate,
		Sex:               r.Sex,
		Height:            r.Height,
		LocationID:        r.LocationID,
		NutritionalStatus: r.NutritionalStatus,
	}
}

// ToUpdateModel mengonversi UpdateToddlerRequest menjadi models.Toddler (hanya field yang dikirim)
func (r *UpdateToddlerRequest) ToUpdateModel(userID int, url string, parentID int) models.Toddler {
	toddlerMapping := models.Toddler{
		UpdatedByID: userID,
	}

	if r.Name != nil {
		toddlerMapping.Name = *r.Name
	}
	if r.Birthdate != nil {
		toddlerMapping.Birthdate = *r.Birthdate
	}
	if r.Sex != "" {
		toddlerMapping.Sex = r.Sex
	}
	if r.Height != nil {
		toddlerMapping.Height = *r.Height
	}
	if r.PhoneNumber != nil {
		toddlerMapping.ParentID = parentID
	}
	if url != "" {
		toddlerMapping.ProfilePicture = url
	}
	if r.LocationID != nil {
		toddlerMapping.LocationID = *r.LocationID
	}
	if r.NutritionalStatus != nil {
		toddlerMapping.NutritionalStatus = *r.NutritionalStatus
	}

	return toddlerMapping
}

// ToCreateRequest mengonversi UpdateToddlerRequest menjadi CreateToddlerRequest
// Berguna untuk memanggil PredictService di dalam fungsi Update
func (r *UpdateToddlerRequest) ToCreateRequest() CreateToddlerRequest {
	req := CreateToddlerRequest{}
	if r.Name != nil {
		req.Name = *r.Name
	}
	if r.Birthdate != nil {
		req.Birthdate = *r.Birthdate
	}
	if r.Sex != "" {
		req.Sex = r.Sex
	}
	if r.Height != nil {
		req.Height = *r.Height
	}
	if r.NutritionalStatus != nil {
		req.NutritionalStatus = *r.NutritionalStatus
	}
	if r.LocationID != nil {
		req.LocationID = *r.LocationID
	}
	if r.PhoneNumber != nil {
		req.PhoneNumber = *r.PhoneNumber
	}
	return req
}
