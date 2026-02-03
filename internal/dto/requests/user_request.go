package requests

import (
	"grovia/internal/models"
	"mime/multipart"
)

type CreateUserRequest struct {
	Name           string                `form:"name" validate:"required"`
	PhoneNumber    string                `form:"phoneNumber" validate:"required,phone"`
	Address        string                `form:"address" validate:"required"`
	Nik            string                `form:"nik" validate:"required,nik"`
	Role           string                `form:"role" validate:"required,oneof=admin kepala_posyandu kader"`
	Password       string                `form:"password" validate:"required,min=6"`
	ProfilePicture *multipart.FileHeader `form:"profilePicture"`
	LocationID     int                   `form:"locationID"`
}

type UpdateUserRequest struct {
	Name           *string               `form:"name,omitempty" validate:"omitempty"`
	PhoneNumber    *string               `form:"phoneNumber,omitempty" validate:"omitempty,phone"`
	Address        *string               `form:"address,omitempty" validate:"omitempty"`
	Nik            *string               `form:"nik,omitempty" validate:"omitempty,nik"`
	Role           *string               `form:"role,omitempty" validate:"omitempty,oneof=admin kepala_posyandu kader"`
	Password       *string               `form:"password,omitempty" validate:"omitempty,min=6"`
	ProfilePicture *multipart.FileHeader `form:"profilePicture,omitempty"`
	LocationID     *int                  `form:"locationID,omitempty"`
}

func (r *CreateUserRequest) ToModel(hashedPassword string, profileURL string, createdBy string) models.User {
	return models.User{
		Name:           r.Name,
		PhoneNumber:    r.PhoneNumber,
		Address:        r.Address,
		Nik:            r.Nik,
		Role:           r.Role,
		Password:       hashedPassword,
		ProfilePicture: profileURL,
		LocationID:     r.LocationID,
		CreatedBy:      createdBy,
		IsActive:       true,
	}
}

func (r *UpdateUserRequest) ToUpdateModel(url string, hashedPassword string) models.User {
	userMapping := models.User{}

	if r.Name != nil {
		userMapping.Name = *r.Name
	}
	if r.PhoneNumber != nil {
		userMapping.PhoneNumber = *r.PhoneNumber
	}
	if r.Address != nil {
		userMapping.Address = *r.Address
	}
	if r.Nik != nil {
		userMapping.Nik = *r.Nik
	}
	if r.LocationID != nil {
		userMapping.LocationID = *r.LocationID
	}
	if r.Role != nil {
		userMapping.Role = *r.Role
	}

	if url != "" {
		userMapping.ProfilePicture = url
	}
	if hashedPassword != "" {
		userMapping.Password = hashedPassword
	}

	return userMapping
}