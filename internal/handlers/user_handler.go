package handlers

import (
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) CreateUser(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	role := ctx.Locals("role")
	locationID := ctx.Locals("location_id").(*int)

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

	var req requests.CreateUserRequest

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

	user, err := u.service.CreateUser(req, role.(string), locationID)

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
		Message: "Create User Success",
		Data:    user,
		Error:   nil,
	})
}

func (u *UserHandler) GetCurrentUser(ctx *fiber.Ctx) error {
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

	user, err := u.service.GetCurrentUser(userID)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Internal Server Error",
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Get User Success",
		Data:    user,
		Error:   nil,
	})
}

func (u *UserHandler) GetUserByID(ctx *fiber.Ctx) error {
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

	user, err := u.service.GetUserById(userID, role.(string))

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Internal Server Error",
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Get User Success",
		Data:    user,
		Error:   nil,
	})
}

func (u *UserHandler) GetUsersByRole(ctx *fiber.Ctx) error {
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

	user, err := u.service.GetUsersByRole(role.(string))

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Internal Server Error",
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Get User Success",
		Data:    user,
		Error:   nil,
	})
}

func (u *UserHandler) UpdateCurrentUser(ctx *fiber.Ctx) error {
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

	var req requests.UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid Request",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: "Invalid Request",
			},
		})
	}

	user, err := u.service.UpdateCurrentUser(userID, req)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Internal Server Error",
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Update User Success",
		Data:    user,
		Error:   nil,
	})
}

func (u *UserHandler) UpdateUserByID(ctx *fiber.Ctx) error {
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

	var req requests.UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid Request",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: "Invalid Request",
			},
		})
	}

	user, err := u.service.UpdateUserByID(userID, req, role.(string))

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Internal Server Error",
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Update User Success",
		Data:    user,
		Error:   nil,
	})
}

func (u *UserHandler) DeleteCurrentUser(ctx *fiber.Ctx) error {
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

	if err := u.service.DeleteCurrentUser(userID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Internal Server Error",
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Delete User Success",
		Error:   nil,
	})
}

func (u *UserHandler) DeleteUserByID(ctx *fiber.Ctx) error {
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

	if err := u.service.DeleteUserByID(userID, role.(string)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.BaseResponse{
			Success: false,
			Message: "Internal Server Error",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Internal Server Error",
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Delete User Success",
		Error:   nil,
	})
}
