package configs

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConnections *gorm.DB

func DBInitiator() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("migration.db.postgres.db_host"),
		viper.GetInt("migration.db.postgres.db_port"),
		viper.GetString("migration.db.postgres.db_user"),
		viper.GetString("migration.db.postgres.db_password"),
		viper.GetString("migration.db.postgres.db_name"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	DBConnections = db

	fmt.Println("Successfully connected to database")
}
