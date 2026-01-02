package routes

import (
	"grovia/internal/handlers"
	"grovia/internal/repositories"
	"grovia/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthRouter(db *gorm.DB, app *fiber.App) {
	var (
		authRepo    = repositories.NewAuthRepository(db)
		authService = services.NewAuthService(authRepo)
		authHandler = handlers.NewAuthHandler(authService)
	)

	public := app.Group("/api/auth")

	public.Post("/login", authHandler.Login)
	public.Post("/reset-password", authHandler.ResetPassword)
	public.Post("/refresh-token", authHandler.RefreshToken)
}
