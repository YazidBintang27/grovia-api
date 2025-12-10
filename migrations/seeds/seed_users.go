package seeds

import (
	"fmt"
	"grovia/internal/models"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB, locations []models.Location) []models.User {
    var count int64
    db.Model(&models.User{}).Count(&count)

    if count > 0 {
        var existing []models.User
        db.Find(&existing)
        fmt.Println("ℹ️ Users already exist. Skipping seeding.")
        return existing
    }

    var users []models.User

    for i, loc := range locations[1:] {
        realIndex := i + 1

        users = append(users, models.User{
            Name:        "Kepala Posyandu " + loc.Name,
            PhoneNumber: fmt.Sprintf("08210%04d", realIndex),
            Address:     loc.Address,
            Nik:         fmt.Sprintf("111111111%03d", realIndex),
            Role:        "kepala_posyandu",
            Password:    "password123",
            LocationID:  int(loc.ID),
        })

        for j := 1; j <= 3; j++ {
            users = append(users, models.User{
                Name:        fmt.Sprintf("Kader %d - %s", j, loc.Name),
                PhoneNumber: fmt.Sprintf("08222%04d", realIndex*3+j),
                Address:     loc.Address,
                Nik:         fmt.Sprintf("222222222%03d", realIndex*3+j),
                Role:        "kader",
                Password:    "password123",
                LocationID:  int(loc.ID),
            })
        }
    }

    db.Create(&users)

    fmt.Println("✅ Seeded users for locations starting from Location ID 2.")
    return users
}

