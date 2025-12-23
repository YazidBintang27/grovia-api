package handlers

import (
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PredictHandler struct {
	service services.PredictService
}

func NewPredictHandler(service services.PredictService) *PredictHandler {
	return &PredictHandler{service: service}
}

func (h *PredictHandler) CreateGroupPredict(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)

	if !ok || userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(responses.BaseResponse{
			Success: false,
			Message: "Unauthorized",
			Error: responses.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Unauthorized",
			},
		})
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "File is required",
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	filePath := "./uploads/" + file.Filename
	if err := ctx.SaveFile(file, filePath); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Failed to save file",
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
	}

	data, err := h.service.CreateGroupPredict(filePath)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Failed to process group predict",
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
	}

	ctx.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Set("Content-Disposition", "attachment; filename=hasil_prediksi.xlsx")
	return ctx.Send(data)
}

func (h *PredictHandler) GetAllPredict(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	locationID := ctx.Locals("location_id").(int)
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")

	if !ok || userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(responses.BaseResponse{
			Success: false,
			Message: "Unauthorized",
			Error: responses.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Unauthorized",
			},
		})
	}

	predicts, meta, err := h.service.GetAllPredict(locationID, pageStr, limitStr)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Get all predict success",
		Data:    predicts,
		Meta:    meta,
	})
}

func (h *PredictHandler) GetAllPredictByToddlerID(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	locationID := ctx.Locals("location_id").(int)

	if !ok || userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(responses.BaseResponse{
			Success: false,
			Message: "Unauthorized",
			Error: responses.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Unauthorized",
			},
		})
	}

	toddlerIDParam := ctx.Params("id")
	toddlerID, err := strconv.Atoi(toddlerIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid toddler ID",
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	predicts, err := h.service.GetAllPredictByToddlerID(locationID, toddlerID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Get all predict by toddler ID success",
		Data:    predicts,
	})
}

func (h *PredictHandler) GetPredictByID(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)

	if !ok || userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(responses.BaseResponse{
			Success: false,
			Message: "Unauthorized",
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
			Message: "Invalid predict ID",
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	predict, err := h.service.GetPredictByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Get predict success",
		Data:    predict,
	})
}

func (h *PredictHandler) UpdatePredictByID(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)

	if !ok || userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(responses.BaseResponse{
			Success: false,
			Message: "Unauthorized",
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
			Message: "Invalid predict ID",
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	var req requests.UpdatePredictRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid Request",
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	updated, err := h.service.UpdatePredictByID(id, &req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Update predict success",
		Data:    updated,
	})
}

func (h *PredictHandler) DeletePredictByID(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	locationID := ctx.Locals("location_id").(int)

	if !ok || userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(responses.BaseResponse{
			Success: false,
			Message: "Unauthorized",
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
			Message: "Invalid predict ID",
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	if err := h.service.DeletePredictByID(id, locationID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Delete predict success",
	})
}

func (h *PredictHandler) GetAllPredictAllLocation(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
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

	predictResponses, meta, err := h.service.GetAllPredictAllLocation(pageStr, limitStr)

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
		Message: "Get All Predict Data Without Location Success",
		Data:    predictResponses,
		Meta:    meta,
		Error:   nil,
	})
}
