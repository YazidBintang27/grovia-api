package requests

import "mime/multipart"

type CreateUserRequest struct {
	Name           string                `form:"name" validate:"required"`
	PhoneNumber    string                `form:"phoneNumber" validate:"required"`
	Address        string                `form:"address" validate:"required"`
	Nik            string                `form:"nik" validate:"required"`
	Role           string                `form:"role" validate:"required"`
	Password       string                `form:"password" validate:"required"`
	ProfilePicture *multipart.FileHeader `form:"profilePicture"`
	LocationID     *int                  `form:"locationId"`
}

type UpdateUserRequest struct {
	Name           string                `form:"name" validate:"required"`
	PhoneNumber    string                `form:"phoneNumber" validate:"required"`
	Address        string                `form:"address" validate:"required"`
	Nik            string                `form:"nik" validate:"required"`
	Role           string                `form:"role" validate:"required"`
	Password       string                `form:"password" validate:"required"`
	ProfilePicture *multipart.FileHeader `form:"profilePicture"`
	LocationID     *int                  `form:"locationId"` 
}
