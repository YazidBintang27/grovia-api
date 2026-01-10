package pkg

import (
	"errors"
	"grovia/internal/dto/responses"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func HandleServiceError(ctx *fiber.Ctx, err error) error {
	var customErr *CustomError

	if errors.As(err, &customErr) {
		return ctx.Status(customErr.StatusCode).JSON(responses.BaseResponse{
			Success: false,
			Message: customErr.Message,
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    customErr.Code,
				Message: customErr.Message,
			},
		})
	}

	return ctx.Status(http.StatusInternalServerError).JSON(responses.BaseResponse{
		Success: false,
		Message: "Internal Server Error",
		Data:    nil,
		Error: responses.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		},
	})
}
