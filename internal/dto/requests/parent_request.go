package requests

type CreateParentRequest struct {
	Name        string `json:"name" validate:"required"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Nik         string `json:"nik" validate:"required"`
	Job         string `json:"job" validate:"required"`
	LocationID  int    `json:"locationID" validate:"required"`
}

type UpdateParentRequest struct {
	Name        *string `json:"name,omitempty"`
	Address     *string `json:"address,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`
	Nik         *string `json:"nik,omitempty"`
	Job         *string `json:"job,omitempty"`
	LocationID  *int    `json:"locationID,omitempty"`
}
