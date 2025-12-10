package seeds

import (
	"fmt"
	"grovia/internal/models"
	"strconv"

	"gorm.io/gorm"
)

func SeedParents(db *gorm.DB, locations []models.Location) []models.Parent {
	var count int64
	db.Model(&models.Parent{}).Count(&count)

	if count > 0 {
		var existing []models.Parent
		db.Find(&existing)
		fmt.Println("ℹ️ Parents already exist. Skipping seeding.")
		return existing
	}

	names := []string{
		"Andi", "Budi", "Citra", "Dewi",
		"Eko", "Farah", "Gilang", "Hana",
		"Iwan", "Joko", "Kirana", "Lia",
		"Made", "Nina", "Oki", "Putri",
		"Qori", "Rina", "Sakti", "Tari",
	}

	var parents []models.Parent

	for i := 0; i < 20; i++ {
		locIndex := (i / 4) + 1
		
		if locIndex >= len(locations) {
			locIndex = len(locations) - 1
		}
		
		loc := locations[locIndex]

		parents = append(parents, models.Parent{
			Name:        names[i],
			Address:     "Alamat RT 0" + strconv.Itoa((i%5)+1) + " / RW 05",
			PhoneNumber: fmt.Sprintf("08123%04d", i+1),
			Nik:         fmt.Sprintf("321654%06d", i+1),
			Job:         "Ibu Rumah Tangga",
			LocationID:  int(loc.ID),
		})
	}

	db.Create(&parents)
	fmt.Println("✅ Seeded 20 parents.")
	return parents
}
