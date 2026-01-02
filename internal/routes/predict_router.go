package routes

import (
	"grovia/internal/handlers"
	"grovia/internal/middlewares"
	"grovia/internal/repositories"
	"grovia/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PredictRouter(db *gorm.DB, app *fiber.App, mlAPIURL string) {
	var (
		predictRepo    = repositories.NewPredictRepository(db)
		predictService = services.NewPredictService(predictRepo, mlAPIURL)
		predictHandler = handlers.NewPredictHandler(predictService)
	)

	r := app.Group("/api/predicts")

	r.Use(middlewares.JWTAuth())

	r.Post("/group", predictHandler.CreateGroupPredict)

	r.Get("/", predictHandler.GetAllPredict)

	r.Get("/all", predictHandler.GetAllPredictAllLocation)

	r.Get("/toddler/:id", predictHandler.GetAllPredictByToddlerID)

	r.Get("/:id", predictHandler.GetPredictByID)

	r.Patch("/:id", predictHandler.UpdatePredictByID)

	r.Delete("/:id", predictHandler.DeletePredictByID)
}
