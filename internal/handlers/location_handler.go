package handlers

import (
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/services"
	"grovia/pkg"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type LocationHandler struct {
	service services.LocationService
}

func NewLocationHandler(service services.LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

// CreateLocation godoc
// @Summary Create new location
// @Description Create a new location (Admin only)
// @Tags Locations
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Location name"
// @Param address formData string true "Address"
// @Param picture formData file false "Location picture"
// @Success 200 {object} responses.BaseResponse
// @Failure 400 {object} responses.BaseResponse
// @Failure 401 {object} responses.BaseResponse
// @Failure 403 {object} responses.BaseResponse
// @Router /locations [post]
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
		return pkg.HandleServiceError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Create location Data Success",
		Data:    locationResponse,
		Error:   nil,
	})
}

// GetAllLocation godoc
// @Summary Get all locations
// @Description Get list of all locations (Admin only)
// @Tags Locations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Search by name"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} responses.BaseResponse
// @Failure 401 {object} responses.BaseResponse
// @Failure 403 {object} responses.BaseResponse
// @Router /locations [get]
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
		return pkg.HandleServiceError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Get All location Data Success",
		Data:    locationsResponse,
		Meta:    meta,
		Error:   nil,
	})
}

// GetLocationByID godoc
// @Summary Get location by ID
// @Description Get location details by ID
// @Tags Locations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Location ID"
// @Success 200 {object} responses.BaseResponse
// @Failure 400 {object} responses.BaseResponse
// @Failure 401 {object} responses.BaseResponse
// @Failure 404 {object} responses.BaseResponse
// @Router /locations/{id} [get]
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
		return pkg.HandleServiceError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Get location Data Success",
		Data:    locationResponse,
		Error:   nil,
	})
}

// UpdateLocationByID godoc
// @Summary Update location by ID
// @Description Update location data by ID (Admin only)
// @Tags Locations
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Location ID"
// @Param name formData string false "Location name"
// @Param address formData string false "Address"
// @Param picture formData file false "Location picture"
// @Success 200 {object} responses.BaseResponse
// @Failure 400 {object} responses.BaseResponse
// @Failure 401 {object} responses.BaseResponse
// @Failure 403 {object} responses.BaseResponse
// @Failure 404 {object} responses.BaseResponse
// @Router /locations/{id} [patch]
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
		return pkg.HandleServiceError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Update location Data Success",
		Data:    locationResponse,
		Error:   nil,
	})
}

// DeleteLocationByID godoc
// @Summary Delete location by ID
// @Description Delete location by ID (Admin only)
// @Tags Locations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Location ID"
// @Success 200 {object} responses.BaseResponse
// @Failure 400 {object} responses.BaseResponse
// @Failure 401 {object} responses.BaseResponse
// @Failure 403 {object} responses.BaseResponse
// @Failure 404 {object} responses.BaseResponse
// @Router /locations/{id} [delete]
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
		return pkg.HandleServiceError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Delete location Data Success",
		Data:    nil,
		Error:   nil,
	})
}
