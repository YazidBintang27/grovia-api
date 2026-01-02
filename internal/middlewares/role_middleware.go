package middlewares

import (
	"grovia/internal/dto/responses"

	"github.com/gofiber/fiber/v2"
)

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusForbidden).JSON(responses.BaseResponse{
				Success: false,
				Message: "Access Denied",
				Data:    nil,
				Error: responses.ErrorResponse{
					Code:    "ACCESS_DENIED",
					Message: "No role information",
				},
			})
		}

		userRole, ok := role.(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
				Success: false,
				Message: "Invalid role type",
				Data:    nil,
				Error: responses.ErrorResponse{
					Code:    "INVALID_ROLE_TYPE",
					Message: "Role should be string",
				},
			})
		}

		for _, allowed := range allowedRoles {
			if userRole == allowed {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(responses.BaseResponse{
			Success: false,
			Message: "Insufficient Permissions",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INSUFFICIENT_PERMISSIONS",
				Message: "You are not allowed to access this resource",
			},
		})
	}
}
