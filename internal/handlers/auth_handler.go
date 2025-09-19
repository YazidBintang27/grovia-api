package handlers

import (
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (a *AuthHandler) ResetPassword(ctx *fiber.Ctx) error {
	var req requests.ResetPasswordRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid request",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	err := a.service.ResetPassword(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Reset password failed",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "RESET_PASSWORD_FAILED",
				Message: err.Error(),
			},
		})
	}

	return ctx.JSON(responses.BaseResponse{
		Success: true,
		Message: "Reset password success",
		Data:    req.PhoneNumber,
		Error:   nil,
	})
}

func (a *AuthHandler) Login(ctx *fiber.Ctx) error {
	var req requests.LoginRequest

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

	token, err := a.service.Login(req)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid Credentials",
			Data:    nil,
			Error: responses.ErrorResponse{
				Code:    "INVALID_CREDENTIALS",
				Message: "Invalid Credentials",
			},
		})
	}

	loginResponse := responses.LoginResponse{
		Token: responses.TokenResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Login Success",
		Data:    loginResponse,
		Error:   nil,
	})
}
