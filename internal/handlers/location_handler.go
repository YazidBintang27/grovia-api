package handlers

import (
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type LocationHandler struct {
	service services.LocationService
}

func NewLocationHandler(service services.LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

func (l *LocationHandler) CreateLocation(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	role := ctx.Locals("role")

	if !ok || userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(responses.BaseResponse{
			Success: false,
			Message: "Unauthorized",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Unauthorized",
			},
		})
	}

	if role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(responses.BaseResponse{
			Success: false,
			Message: "Forbidden",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "FORBIDDEN",
				Message: "Forbidden",
			},
		})
	}

	var req requests.LocationRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid Request",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	file, err := ctx.FormFile("picture")
	if err == nil {
		req.Picture = file
	}

	locationResponse, err := l.service.CreateLocation(ctx.Context(), req, userID)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Create location Data Success",
		Data:    locationResponse,
		Error:   nil,
	})
}

func (l *LocationHandler) GetAllLocation(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	role := ctx.Locals("role")

	name := ctx.Query("name")
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")

	if !ok || userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(responses.BaseResponse{
			Success: false,
			Message: "Unauthorized",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Unauthorized",
			},
		})
	}

	if role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(responses.BaseResponse{
			Success: false,
			Message: "Forbidden",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "FORBIDDEN",
				Message: "Forbidden",
			},
		})
	}

	locationsResponse, meta, err := l.service.GetAllLocation(name, pageStr, limitStr)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Get All location Data Success",
		Data:    locationsResponse,
		Meta:    meta,
		Error:   nil,
	})
}

func (l *LocationHandler) GetLocationByID(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)

	if !ok || userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(responses.BaseResponse{
			Success: false,
			Message: "Unauthorized",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Unauthorized",
			},
		})
	}

	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid Request",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	locationResponse, err := l.service.GetLocationByID(id)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Get location Data Success",
		Data:    locationResponse,
		Error:   nil,
	})
}

func (l *LocationHandler) UpdateLocationByID(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	role := ctx.Locals("role")

	if !ok || userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(responses.BaseResponse{
			Success: false,
			Message: "Unauthorized",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Unauthorized",
			},
		})
	}

	if role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(responses.BaseResponse{
			Success: false,
			Message: "Forbidden",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "FORBIDDEN",
				Message: "Forbidden",
			},
		})
	}

	var req requests.LocationRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid Request",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	file, err := ctx.FormFile("picture")
	if err == nil {
		req.Picture = file
	}

	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid Request",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	locationResponse, err := l.service.UpdateLocationByID(ctx.Context(), id, userID, req)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Update location Data Success",
		Data:    locationResponse,
		Error:   nil,
	})
}

func (l *LocationHandler) DeleteLocationByID(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	role := ctx.Locals("role")

	if !ok || userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(responses.BaseResponse{
			Success: false,
			Message: "Unauthorized",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Unauthorized",
			},
		})
	}

	if role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(responses.BaseResponse{
			Success: false,
			Message: "Forbidden",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "FORBIDDEN",
				Message: "Forbidden",
			},
		})
	}

	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid Request",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	err = l.service.DeleteLocationByID(id, userID)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Delete location Data Success",
		Data:    nil,
		Error:   nil,
	})
}
