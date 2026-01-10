package pkg

import "net/http"

type CustomError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) *CustomError {
	return &CustomError{
		StatusCode: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    message,
	}
}

func NewUnauthorizedError(message string) *CustomError {
	return &CustomError{
		StatusCode: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    message,
	}
}

func NewForbiddenError(message string) *CustomError {
	return &CustomError{
		StatusCode: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    message,
	}
}

func NewNotFoundError(message string) *CustomError {
	return &CustomError{
		StatusCode: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    message,
	}
}

func NewInternalServerError(message string) *CustomError {
	return &CustomError{
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    message,
	}
}

func NewConflictError(message string) *CustomError {
	return &CustomError{
		StatusCode: http.StatusConflict,
		Code:       "CONFLICT",
		Message:    message,
	}
}

func NewUnprocessableEntityError(message string) *CustomError {
	return &CustomError{
		StatusCode: http.StatusUnprocessableEntity,
		Code:       "UNPROCESSABLE_ENTITY",
		Message:    message,
	}
}