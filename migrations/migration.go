package migrations

import (
	"embed"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/gorm"
)

//go:embed sql_migrations/*.sql
var dbMigrations embed.FS

func Initiator(dbParam *gorm.DB) {
	var dbName string
	dbParam.Raw("SELECT current_database()").Scan(&dbName)
	fmt.Println("MIGRATION CONNECTED TO DB:", dbName)
	gormDB, err := dbParam.DB()
	if err != nil {
		panic("failed to get gorm.DB: " + err.Error())
	}

	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "sql_migrations",
	}

	n, errs := migrate.Exec(gormDB, "postgres", migrations, migrate.Up)
	if errs != nil {
		panic(errs)
	}

	fmt.Println("Migrations success, applied", n, "migrations!")
}
