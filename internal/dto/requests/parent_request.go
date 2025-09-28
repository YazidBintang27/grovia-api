package requests

type CreateParentWithToddlersRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Address     string                 `json:"address" validate:"required"`
	PhoneNumber string                 `json:"phone_number" validate:"required"`
	Nik         string                 `json:"nik" validate:"required"`
	Job         string                 `json:"job" validate:"required"`
	LocationID  *int                   `json:"location_id" validate:"required"`
	Toddlers    []CreateToddlerRequest `json:"toddlers" validate:"required,dive"`
}

type UpdateParentRequest struct {
	Name        string `json:"name" validate:"required"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Nik         string `json:"nik" validate:"required"`
	Job         string `json:"job" validate:"required"`
	LocationID  *int   `json:"location_id" validate:"required"`
}
