package requests

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
