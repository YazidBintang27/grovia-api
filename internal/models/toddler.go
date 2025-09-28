package models

import "time"

type Toddler struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ParentID       int       `json:"parent_id" gorm:"not null"`
	Parent         Parent    `json:"parent" gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	LocationID     int       `json:"location_id" gorm:"not null"`
	Location       Location  `json:"location" gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name           string    `json:"name" gorm:"type:varchar(100);not null"`
	Birthdate      time.Time `json:"birthdate" gorm:"type:date;not null"`
	Gender         string    `json:"gender" gorm:"type:varchar(10);not null"`
	Height         float64   `json:"height" gorm:"type:decimal(5,2)"`
	ProfilePicture string    `json:"profile_picture" gorm:"type:varchar(100)"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
