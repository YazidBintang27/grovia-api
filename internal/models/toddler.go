package models

import "time"

type Toddler struct {
	ID                int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ParentID          int       `json:"parentId" gorm:"not null"`
	Parent            Parent    `json:"parent" gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	LocationID        int       `json:"locationId" gorm:"not null"`
	Location          Location  `json:"location" gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name              string    `json:"name" gorm:"type:varchar(100);not null"`
	Birthdate         time.Time `json:"birthdate" gorm:"type:date;not null"`
	Sex               string    `json:"sex" gorm:"type:varchar(10);not null"`
	Height            float64   `json:"height" gorm:"type:decimal(4,1)"`
	ProfilePicture    string    `json:"profilePicture" gorm:"type:text"`
	NutritionalStatus string    `json:"nutritionalStatus" gorm:"type:varchar(50)"`
	CreatedAt         time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
