package seeds

import (
	"fmt"
	"grovia/internal/models"
	"time"

	"gorm.io/gorm"
)

func SeedToddlers(db *gorm.DB, locations []models.Location, parents []models.Parent) []models.Toddler {
	var count int64
	db.Model(&models.Toddler{}).Count(&count)

	if count > 0 {
		var existing []models.Toddler
		db.Find(&existing)
		fmt.Println("ℹ️ Toddlers already exist. Skipping seeding.")
		return existing
	}

	names := []string{
		"Rafi", "Salsa", "Tono", "Umi",
		"Vira", "Wawan", "Xena", "Yogi",
		"Zara", "Adit", "Bella", "Ciko",
		"Dara", "Evan", "Fina", "Gio",
		"Hani", "Ivan", "Jeni", "Koko",
	}

	var toddlers []models.Toddler

	for i := 0; i < 20; i++ {

		locIndex := (i / 4) + 1

		if locIndex >= len(locations) {
			locIndex = len(locations) - 1
		}

		loc := locations[locIndex]

		parent := parents[i]

		t := models.Toddler{
			Name:              names[i],
			Sex:               "M",
			Birthdate:         time.Now().AddDate(-3, 0, -i),
			Height:            float64(60 + i),
			NutritionalStatus: "normal",
			LocationID:        int(loc.ID),

			ParentID: parent.ID,
		}

		db.Create(&t)
		toddlers = append(toddlers, t)

		SeedPredictsForToddler(db, t.ID)
	}

	fmt.Println("✅ Seeded 20 toddlers with parent relation + predict history.")
	return toddlers
}
