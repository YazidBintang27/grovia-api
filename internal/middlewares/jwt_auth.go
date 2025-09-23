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

		tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)

		claims, err := pkg.ValidateToken(tokenStr)

		if err != nil {
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

		ctx.Locals("user_id", claims.UserID)
		ctx.Locals("role", claims.Role)

		return ctx.Next()
	}
}
