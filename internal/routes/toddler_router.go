package routes

import (
	"grovia/internal/handlers"
	"grovia/internal/middlewares"
	"grovia/internal/repositories"
	"grovia/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ToddlerRouter(db *gorm.DB, app *fiber.App, s3 *services.S3Service, predict services.PredictService) {
	var (
		toddlerRepo    = repositories.NewToddlerRepository(db)
		parentRepo     = repositories.NewParentRepository(db)
		toddlerService = services.NewToddlerService(toddlerRepo, parentRepo, s3, predict)
		toddlerHandler = handlers.NewToddlerHandler(toddlerService)
	)

	r := app.Group("/api/toddlers")

	r.Use(middlewares.JWTAuth())

	r.Post("/", toddlerHandler.CreateToddler)

	r.Post("/with-parent", toddlerHandler.CreateToddlerWithParent)

	r.Get("/check-toddler", toddlerHandler.CheckToddlerExists)

	r.Get("/", toddlerHandler.GetAllToddler)

	r.Get("/all", toddlerHandler.GetAllToddlerAllLocation)

	r.Patch("/without-predict/:id", toddlerHandler.UpdateToddlerByIDWithoutPredict)

	r.Get("/:id", toddlerHandler.GetToddlerByID)

	r.Patch("/:id", toddlerHandler.UpdateToddlerByID)

	r.Delete("/:id", toddlerHandler.DeleteToddlerByID)

}
