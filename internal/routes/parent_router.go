package routes

import (
	"grovia/internal/handlers"
	"grovia/internal/middlewares"
	"grovia/internal/repositories"
	"grovia/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ParentRouter(db *gorm.DB, app *fiber.App) {
	var (
		parentRepo    = repositories.NewParentRepository(db)
		parentService = services.NewParentService(parentRepo)
		parentHandler = handlers.NewParentHandler(parentService)
	)

	r := app.Group("/api/parents")

	r.Use(middlewares.JWTAuth())

	r.Get("/check-phone", parentHandler.CheckPhoneExists)

	r.Post("/", parentHandler.CreateParent)

	r.Get("/", parentHandler.GetAllParent)

	r.Get("/all", parentHandler.GetAllPredictAllLocation)

	r.Get("/:id", parentHandler.GetParentByID)

	r.Patch("/:id", parentHandler.UpdateParentByID)

	r.Delete("/:id", parentHandler.DeleteParentByID)
}
