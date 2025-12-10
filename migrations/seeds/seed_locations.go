package seeds

import (
	"fmt"
	"grovia/internal/models"

	"gorm.io/gorm"
)

func SeedLocations(db *gorm.DB) []models.Location {
	var count int64
	db.Model(&models.Location{}).Count(&count)

	if count > 0 {
		var existing []models.Location
		db.Find(&existing)
		fmt.Println("ℹ️ Locations already exist. Skipping seeding.")
		return existing
	}

	locations := []models.Location{
		{Name: "Desa Pangadegan", Address: "RT 01 / RW 01"},
		{Name: "Nusa Indah 1", Address: "RT 01 / RW 05"},
		{Name: "Nusa Indah 2", Address: "RT 02 / RW 05"},
		{Name: "Nusa Indah 3", Address: "RT 03 / RW 05"},
		{Name: "Nusa Indah 4", Address: "RT 04 / RW 06"},
		{Name: "Nusa Indah 5", Address: "RT 05 / RW 06"},
	}

	db.Create(&locations)
	fmt.Println("✅ Seeded locations: Desa Pangadegan + Nusa Indah 1–5.")

	return locations
}
