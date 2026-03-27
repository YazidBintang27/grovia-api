package main

import (
	"grovia/configs"
	"grovia/internal/firebase"
	"grovia/internal/repositories"
	"grovia/internal/routes"
	"grovia/internal/services"
	"grovia/migrations"

	// "grovia/migrations"
	"grovia/migrations/seeds"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "grovia/docs"
)

// @title Grovia API
// @version 1.0
// @description API documentation for Grovia App
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@grovia.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	firebase.InitFirebase()

	cfg := configs.LoadConfig()

	configs.DBInitiator()
	db := configs.DBConnections

	migrations.Initiator(db)
	// seeds.SeedDefaultLocation(db)
	seeds.SeedAdmin(db)
	// RunAllSeeds(db)

	s3 := services.NewS3Service(cfg.Aws)
	predictRepo := repositories.NewPredictRepository(db)
	predictService := services.NewPredictService(predictRepo, cfg.MLAPIURL)

	InitiateRoutes(db, s3, predictService, cfg.MLAPIURL)
}

func InitiateRoutes(db *gorm.DB, s3 *services.S3Service, predict services.PredictService, mlAPIURL string) {
	app := fiber.New()

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	routes.AuthRouter(db, app)
	routes.LocationRouter(app, db, s3)
	routes.ParentRouter(db, app)
	routes.PredictRouter(db, app, mlAPIURL)
	routes.ToddlerRouter(db, app, s3, predict)
	routes.UserRouter(app, db, s3)

	log.Fatal(app.Listen(":8080"))
}

func RunAllSeeds(db *gorm.DB) {
	locations := seeds.SeedLocations(db)
	parents := seeds.SeedParents(db, locations)
	seeds.SeedToddlers(db, locations, parents)
	seeds.SeedUsers(db, locations)
}
