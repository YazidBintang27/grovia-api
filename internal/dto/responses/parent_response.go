package responses

import "time"

type ParentResponse struct {
	ID          int       `json:"id"`
	LocationID  int       `json:"location_id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	Nik         string    `json:"nik"`
	Job         string    `json:"job"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ParentWithToddlersResponse struct {
	ID          int               `json:"id"`
	LocationID  int               `json:"location_id"`
	Name        string            `json:"name"`
	PhoneNumber string            `json:"phone_number"`
	Address     string            `json:"address"`
	Nik         string            `json:"nik"`
	Job         string            `json:"job"`
	Toddlers    []ToddlerResponse `json:"toddlers"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
