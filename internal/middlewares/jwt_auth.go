package middlewares

import (
	"grovia/internal/dto/responses"
	"grovia/pkg"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTAuth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(responses.BaseResponse{
				Success: false,
				Message: "Authorization header missing",
				Data:    nil,
				Error: responses.ErrorResponse{
					Code:    "UNAUTHORIZED",
					Message: "Authorization header missing",
				},
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		tokenStr = strings.TrimSpace(tokenStr)

		claims, err := pkg.ValidateToken(tokenStr)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(responses.BaseResponse{
				Success: false,
				Message: "Invalid or expired token",
				Data:    nil,
				Error: responses.ErrorResponse{
					Code:    "INVALID_TOKEN",
					Message: err.Error(),
				},
			})
		}

		ctx.Locals("user_id", claims.UserID)
		ctx.Locals("role", claims.Role)
		ctx.Locals("location_id", claims.LocationID)

		return ctx.Next()
	}
}
