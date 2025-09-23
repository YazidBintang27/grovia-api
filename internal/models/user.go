package models

import "time"

type User struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	LocationID     int       `json:"location_id" gorm:"not null"` 
	Location       Location  `json:"location" gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Name           string    `json:"name" gorm:"type:varchar(100);not null"`
	PhoneNumber    string    `json:"phone_number" gorm:"type:varchar(20);unique;not null"`
	Address        string    `json:"address" gorm:"type:text"`
	Nik            string    `json:"nik" gorm:"type:varchar(16);unique;not null"`
	ProfilePicture string    `json:"profile_picture" gorm:"type:text"`
	Password       string    `json:"-" gorm:"type:varchar(255);not null"`
	Role           string    `json:"role" gorm:"type:varchar(50);default:'user'"`
	CreatedBy      string    `json:"created_by" gorm:"type:varchar(100)"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
