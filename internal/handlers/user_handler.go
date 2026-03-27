package handlers

import (
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/services"
	"grovia/pkg"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser godoc
// @Summary Create new user
// @Description Create a new user (Admin and Kepala Posyandu only)
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "User name"
// @Param phoneNumber formData string true "Phone number"
// @Param password formData string true "Password"
// @Param role formData string true "User role"
// @Param profilePicture formData file false "Profile picture"
// @Success 200 {object} responses.BaseResponse
// @Failure 400 {object} responses.BaseResponse
// @Failure 401 {object} responses.BaseResponse
// @Failure 403 {object} responses.BaseResponse
// @Router /users [post]
func (u *UserHandler) CreateUser(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	role := ctx.Locals("role").(string)
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

	file, err := ctx.FormFile("profilePicture")
	if err == nil {
		req.ProfilePicture = file
	}

	user, err := u.service.CreateUser(ctx.Context(), req, role, locationID)

	if err != nil {
		return pkg.HandleServiceError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Create User Success",
		Data:    user,
		Error:   nil,
	})
}

// GetCurrentUser godoc
// @Summary Get current user profile
// @Description Get logged in user profile
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.BaseResponse
// @Failure 401 {object} responses.BaseResponse
// @Router /users/current [get]
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
		return pkg.HandleServiceError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Get User Success",
		Data:    user,
		Error:   nil,
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user details by ID (Admin and Kepala Posyandu only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} responses.BaseResponse
// @Failure 400 {object} responses.BaseResponse
// @Failure 401 {object} responses.BaseResponse
// @Failure 403 {object} responses.BaseResponse
// @Failure 404 {object} responses.BaseResponse
// @Router /users/{id} [get]
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

	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid target ID",
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	user, err := u.service.GetUserById(id, role.(string))

	if err != nil {
		return pkg.HandleServiceError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Get User Success",
		Data:    user,
		Error:   nil,
	})
}

// GetUsersByRole godoc
// @Summary Get users by role
// @Description Get list of users filtered by role (Admin and Kepala Posyandu only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Search by name"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} responses.BaseResponse
// @Failure 401 {object} responses.BaseResponse
// @Failure 403 {object} responses.BaseResponse
// @Router /users [get]
func (u *UserHandler) GetUsersByRole(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int)
	locationID := ctx.Locals("location_id").(int)
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

	user, meta, err := u.service.GetUsersByRole(role.(string), name, pageStr, limitStr, locationID)

	if err != nil {
		return pkg.HandleServiceError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Get User Success",
		Data:    user,
		Meta:    meta,
		Error:   nil,
	})
}

// UpdateCurrentUser godoc
// @Summary Update current user profile
// @Description Update logged in user profile
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string false "User name"
// @Param phoneNumber formData string false "Phone number"
// @Param password formData string false "New password"
// @Param profilePicture formData file false "Profile picture"
// @Success 200 {object} responses.BaseResponse
// @Failure 400 {object} responses.BaseResponse
// @Failure 401 {object} responses.BaseResponse
// @Router /users/current [patch]
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
				Message: err.Error(),
			},
		})
	}

	file, err := ctx.FormFile("profilePicture")
	if err == nil {
		req.ProfilePicture = file
	}

	user, err := u.service.UpdateCurrentUser(ctx.Context(), userID, req)

	if err != nil {
		return pkg.HandleServiceError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Update User Success",
		Data:    user,
		Error:   nil,
	})
}

// UpdateUserByID godoc
// @Summary Update user by ID
// @Description Update user data by ID (Admin and Kepala Posyandu only)
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param name formData string false "User name"
// @Param phoneNumber formData string false "Phone number"
// @Param password formData string false "New password"
// @Param profilePicture formData file false "Profile picture"
// @Success 200 {object} responses.BaseResponse
// @Failure 400 {object} responses.BaseResponse
// @Failure 401 {object} responses.BaseResponse
// @Failure 403 {object} responses.BaseResponse
// @Failure 404 {object} responses.BaseResponse
// @Router /users/{id} [patch]
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
				Message: err.Error(),
			},
		})
	}

	file, err := ctx.FormFile("profilePicture")
	if err == nil {
		req.ProfilePicture = file
	}

	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid target ID",
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	user, err := u.service.UpdateUserByID(ctx.Context(), id, req, role.(string))

	if err != nil {
		return pkg.HandleServiceError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Update User Success",
		Data:    user,
		Error:   nil,
	})
}

// DeleteCurrentUser godoc
// @Summary Delete current user
// @Description Delete logged in user account
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.BaseResponse
// @Failure 401 {object} responses.BaseResponse
// @Router /users/current [delete]
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
		return pkg.HandleServiceError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Delete User Success",
		Error:   nil,
	})
}

// DeleteUserByID godoc
// @Summary Delete user by ID
// @Description Delete user by ID (Admin and Kepala Posyandu only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} responses.BaseResponse
// @Failure 400 {object} responses.BaseResponse
// @Failure 401 {object} responses.BaseResponse
// @Failure 403 {object} responses.BaseResponse
// @Failure 404 {object} responses.BaseResponse
// @Router /users/{id} [delete]
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

	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.BaseResponse{
			Success: false,
			Message: "Invalid target ID",
			Error: responses.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
	}

	if err := u.service.DeleteUserByID(id, role.(string)); err != nil {
		return pkg.HandleServiceError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.BaseResponse{
		Success: true,
		Message: "Delete User Success",
		Error:   nil,
	})
}
