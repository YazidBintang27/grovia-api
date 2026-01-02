package models

import "time"

type Location struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	Address   string    `json:"address" gorm:"type:varchar(100)"`
	Picture   string    `json:"picture" gorm:"type:text"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
