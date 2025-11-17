package handlers

import (
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ParentHandler struct {
	service services.ParentService
}

func NewParentHandler(service services.ParentService) *ParentHandler {
	return &ParentHandler{service: service}
}

func (p *ParentHandler) CreateParent(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	locationID := ctx.Locals("location_id").(int)

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

	var req requests.CreateParentRequest
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

	req.LocationID = locationID

	parentResp, err := p.service.CreateParent(req)
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

	return ctx.Status(fiber.StatusCreated).JSON(responses.BaseResponse{
		Success: true,
		Message: "Create Toddler Success",
		Data:    parentResp,
		Error:   nil,
	})
}

func (p *ParentHandler) GetAllParent(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	locationID := ctx.Locals("location_id").(int)

	name := ctx.Query("name")

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

	parents, err := p.service.GetAllParent(locationID, name)

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
		Message: "Get All Parent Data Success",
		Data:    parents,
		Error:   nil,
	})
}

func (p *ParentHandler) GetParentByID(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	locationID := ctx.Locals("location_id").(int)

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

	parent, err := p.service.GetParentByID(id, locationID)

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
		Message: "Get Parent Data Success",
		Data:    parent,
		Error:   nil,
	})
}

func (p *ParentHandler) UpdateParentByID(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	locationID := ctx.Locals("location_id").(int)

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

	var req requests.UpdateParentRequest
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

	parentResponses, err := p.service.UpdateParentByID(id, locationID, req)

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
		Message: "Update Parent Data Success",
		Data:    parentResponses,
		Error:   nil,
	})
}

func (p *ParentHandler) DeleteParentByID(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	locationID := ctx.Locals("location_id").(int)

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

	err = p.service.DeleteParentByID(id, locationID)

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
		Message: "Delete Parent Data Success",
		Data:    nil,
		Error:   nil,
	})
}

func (p *ParentHandler) CheckPhoneExists(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)

	phoneNumber := ctx.Query("phone_number")

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

	if phoneNumber == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid Request",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: "phone_number is required",
			},
		})
	}

	parent, err := p.service.CheckPhoneExists(phoneNumber)
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

	exists := parent != nil

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Check Phone Exists Success",
		Data: fiber.Map{
			"exists": exists,
			"parent_id": func() *int {
				if parent != nil {
					return &parent.ID
				}
				return nil
			}(),
		},
		Error: nil,
	})
}

func (p *ParentHandler) GetAllPredictAllLocation(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)

	name := ctx.Query("name")

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

	parentResponses, err := p.service.GetAllParentAllLocation(name)

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
		Message: "Get All Parent Data Without Location Success",
		Data:    parentResponses,
		Error:   nil,
	})
}
