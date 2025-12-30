package models

import "time"

type Predict struct {
	ID                int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ToddlerID         int       `json:"toddlerId" gorm:"not null"`
	LocationID        int       `json:"locationId" gorm:"not null"`
	CreatedByID       int       `json:"createdByID" gorm:"not null"`
	DeletedByID       *int       `json:"deletedByID"`
	Toddler           Toddler   `json:"toddler" gorm:"foreignKey:ToddlerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Name              string    `json:"name" gorm:"type:varchar(100);not null"`
	Height            float64   `json:"height" gorm:"type:decimal(4,1);not null"`
	Age               int       `json:"age" gorm:"not null"`
	Sex               string    `json:"sex" gorm:"not null"`
	Zscore            float64   `json:"zscore" gorm:"type:decimal(4,1);not null"`
	NutritionalStatus string    `json:"nutritionalStatus" gorm:"type:varchar(50);not null"`
	CreatedAt         time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt         *time.Time `json:"deletedAt" gorm:"index"`
}
