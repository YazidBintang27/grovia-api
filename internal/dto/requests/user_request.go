package requests

import "mime/multipart"

type CreateUserRequest struct {
	Name           string                `form:"name" validate:"required"`
	PhoneNumber    string                `form:"phoneNumber" validate:"required,phone"`
	Address        string                `form:"address" validate:"required"`
	Nik            string                `form:"nik" validate:"required,nik"`
	Role           string                `form:"role" validate:"required,oneof=admin kepala_posyandu kader"`
	Password       string                `form:"password" validate:"required,min6"`
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
