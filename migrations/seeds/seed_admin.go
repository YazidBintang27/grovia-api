package seeds

import (
	"fmt"
	"log"

	"grovia/internal/models"
	"grovia/pkg"

	"gorm.io/gorm"
)
func SeedAdmin(db *gorm.DB) {
	location := SeedDefaultLocation(db)

	var count int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

	if count == 0 {
		hashedPassword, _ := pkg.HashPassword("fufufafa")

		admin := models.User{
			LocationID:     location.ID,
			Name:           "Administrator",
			PhoneNumber:    "08123456789",
			Address:        "Head Quarter",
			Nik:            "0000000000000000",
			Password:       string(hashedPassword),
			Role:           "admin",
			CreatedBy:      "system",
			ProfilePicture: "",
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Println("❌ Failed to create admin user:", err)
		} else {
			fmt.Println("✅ Admin user created successfully (phone: 08123456789, pass: fufufafa)")
		}
	} else {
		fmt.Println("ℹ️ Admin already exists, skipping seed.")
	}
}
