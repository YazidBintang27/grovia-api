package routes

import (
	"grovia/internal/handlers"
	"grovia/internal/middlewares"
	"grovia/internal/repositories"
	"grovia/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LocationRouter(app *fiber.App, db *gorm.DB, s3 *services.S3Service) {
	repo := repositories.NewLocationRepository(db)
	service := services.NewLocationService(repo, s3)
	handler := handlers.NewLocationHandler(service)

	r := app.Group("/api/locations")

	r.Use(middlewares.JWTAuth())

	r.Post("/", middlewares.RoleMiddleware("admin"), handler.CreateLocation)
	r.Get("/", handler.GetAllLocation)
	r.Get("/:id", handler.GetLocationByID)
	r.Patch("/:id", middlewares.RoleMiddleware("admin"), handler.UpdateLocationByID)
	r.Delete("/:id", middlewares.RoleMiddleware("admin"), handler.DeleteLocationByID)
}
