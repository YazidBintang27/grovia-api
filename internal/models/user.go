package models

import "time"

type User struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	LocationID     int       `json:"location_id" gorm:"not null"`
	Location       Location  `json:"location" gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Name           string    `json:"name" gorm:"type:varchar(100);not null"`
	PhoneNumber    string    `json:"phoneNumber" gorm:"type:varchar(100);unique;not null"`
	Address        string    `json:"address" gorm:"type:varchar(100)"`
	Nik            string    `json:"nik" gorm:"type:varchar(100);unique;not null"`
	ProfilePicture string    `json:"profilePicture" gorm:"type:text"`
	Password       string    `json:"-" gorm:"type:varchar(100);not null"`
	Role           string    `json:"role" gorm:"type:varchar(100);"`
	IsActive       bool      `json:"isActive" gorm:"default:true"`
	CreatedBy      string    `json:"createdBy" gorm:"type:varchar(100)"`
	CreatedAt      time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
