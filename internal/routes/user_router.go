package routes

import (
	"grovia/internal/handlers"
	"grovia/internal/middlewares"
	"grovia/internal/repositories"
	"grovia/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRouter(app *fiber.App, db *gorm.DB, s3 *services.S3Service) {
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo, s3)
	handler := handlers.NewUserHandler(service)

	r := app.Group("/api/users")

	r.Use(middlewares.JWTAuth())

	r.Post("/", middlewares.RoleMiddleware("admin", "kepala_posyandu"), handler.CreateUser)

	r.Get("/current", handler.GetCurrentUser)
	r.Patch("/current", handler.UpdateCurrentUser)
	r.Delete("/current", handler.DeleteCurrentUser)

	r.Get("/", middlewares.RoleMiddleware("admin", "kepala_posyandu"), handler.GetUsersByRole)
	r.Get("/:id", middlewares.RoleMiddleware("admin", "kepala_posyandu"), handler.GetUserByID)
	r.Patch("/:id", middlewares.RoleMiddleware("admin", "kepala_posyandu"), handler.UpdateUserByID)
	r.Delete("/:id", middlewares.RoleMiddleware("admin", "kepala_posyandu"), handler.DeleteUserByID)
}
