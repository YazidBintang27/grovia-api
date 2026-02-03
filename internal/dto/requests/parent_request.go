package requests

import (
	"grovia/internal/models"
	"strings"
)

type CreateParentRequest struct {
	Name        string `json:"name" validate:"required"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required,phone"`
	Nik         string `json:"nik" validate:"required,nik"`
	Job         string `json:"job" validate:"required"`
	LocationID  int    `json:"locationID" validate:"required"`
}

type UpdateParentRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty"`
	Address     *string `json:"address,omitempty" validate:"omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty" validate:"omitempty,phone"`
	Nik         *string `json:"nik,omitempty" validate:"omitempty,nik"`
	Job         *string `json:"job,omitempty" validate:"omitempty"`
	LocationID  *int    `json:"locationID,omitempty" validate:"omitempty"`
}

// ToModel mengonversi CreateParentRequest menjadi models.Parent
func (r *CreateParentRequest) ToModel(userID int) models.Parent {
	return models.Parent{
		CreatedByID: userID,
		UpdatedByID: userID,
		Name:        strings.TrimSpace(r.Name),
		Address:     strings.TrimSpace(r.Address),
		PhoneNumber: strings.TrimSpace(r.PhoneNumber),
		Nik:         strings.TrimSpace(r.Nik),
		Job:         strings.TrimSpace(r.Job),
		LocationID:  r.LocationID,
	}
}

// ToUpdateModel mengonversi UpdateParentRequest menjadi models.Parent
// Menangani field opsional (pointer) dan melakukan trimming otomatis
func (r *UpdateParentRequest) ToUpdateModel(userID int) models.Parent {
	parentMapping := models.Parent{
		UpdatedByID: userID,
	}

	if r.Name != nil {
		parentMapping.Name = strings.TrimSpace(*r.Name)
	}
	if r.PhoneNumber != nil {
		parentMapping.PhoneNumber = strings.TrimSpace(*r.PhoneNumber)
	}
	if r.Address != nil {
		parentMapping.Address = strings.TrimSpace(*r.Address)
	}
	if r.Nik != nil {
		parentMapping.Nik = strings.TrimSpace(*r.Nik)
	}
	if r.Job != nil {
		parentMapping.Job = strings.TrimSpace(*r.Job)
	}
	if r.LocationID != nil {
		parentMapping.LocationID = *r.LocationID
	}

	return parentMapping
}
