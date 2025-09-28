package models

import "time"

type Parent struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	LocationID  int       `json:"location_id" gorm:"not null"`
	Location    Location  `json:"location" gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(100);unique;not null"`
	Address     string    `json:"address" gorm:"type:varchar(100)"`
	Nik         string    `json:"nik" gorm:"type:varchar(100);unique;not null"`
	Job         string    `json:"job" gorm:"type:varchar(100)"`
	Toddlers    []Toddler `json:"toddlers" gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
