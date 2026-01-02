package responses

import "time"

type ParentResponse struct {
	ID          int               `json:"id"`
	LocationID  int               `json:"locationID"`
	CreatedByID int               `json:"createdByID"`
	UpdatedByID int               `json:"updatedByID"`
	Name        string            `json:"name"`
	PhoneNumber string            `json:"phoneNumber"`
	Address     string            `json:"address"`
	Nik         string            `json:"nik"`
	Job         string            `json:"job"`
	Toddlers    []ToddlerResponse `json:"toddlers"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}
