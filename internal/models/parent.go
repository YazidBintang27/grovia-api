package models

import "time"

type Parent struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	LocationID  int       `json:"locationId" gorm:"not null"`
	CreatedByID int       `json:"createdByID" gorm:"not null"`
	UpdatedByID int       `json:"updatedByID" gorm:"not null"`
	DeletedByID *int       `json:"deletedByID"`
	Location    Location  `json:"location" gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	PhoneNumber string    `json:"phoneNumber" gorm:"type:varchar(100);unique;not null"`
	Address     string    `json:"address" gorm:"type:varchar(100)"`
	Nik         string    `json:"nik" gorm:"type:varchar(100);unique;not null"`
	Job         string    `json:"job" gorm:"type:varchar(100)"`
	Toddlers    []Toddler `json:"toddlers" gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deletedAt" gorm:"index"`
}
