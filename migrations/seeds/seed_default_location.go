package seeds

import (
	"fmt"
	"log"

	"grovia/internal/models"

	"gorm.io/gorm"
)

func SeedDefaultLocation(db *gorm.DB) models.Location {
	var location models.Location
	result := db.First(&location)

	if result.RowsAffected == 0 {
		defaultLocation := models.Location{
			Name:    "Head Quarter",
			Address: "Default Address",
		}
		if err := db.Create(&defaultLocation).Error; err != nil {
			log.Println("❌ Failed to create default location:", err)
		} else {
			fmt.Println("✅ Default location created.")
		}
		return defaultLocation
	}

	fmt.Println("ℹ️ Using existing location:", location.Name)
	return location
}
