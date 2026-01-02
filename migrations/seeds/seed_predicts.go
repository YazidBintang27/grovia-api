package seeds

import (
	"grovia/internal/models"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func SeedPredictsForToddler(db *gorm.DB, toddlerID int) {
	var toddler models.Toddler
	db.First(&toddler, toddlerID) 

	for i := 0; i < 5; i++ {
		p := models.Predict{
			Name:              toddler.Name,        
			ToddlerID:         toddlerID,
			LocationID:        toddler.LocationID,
			Height:            60 + float64(rand.Intn(10)),
			Age:               12 + i,
			Sex:               toddler.Sex,         
			Zscore:            float64(rand.Intn(200))/100 - 1,
			NutritionalStatus: "normal",
			CreatedAt:         time.Now().AddDate(0, 0, -i),
		}

		db.Create(&p)
	}
}


