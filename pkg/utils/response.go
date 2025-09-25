package utils

import (
	"github.com/gofiber/fiber/v2"
)

// APIResponse represents a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Meta represents pagination metadata
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

func SendSuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SendErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: false,
		Message: message,
	})
}

func SendValidationErrorResponse(c *fiber.Ctx, message string, errors interface{}) error {
	return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
		Success: false,
		Message: message,
		Data:    errors,
	})
}

func SendPaginatedResponse(c *fiber.Ctx, message string, data interface{}, meta *Meta) error {
	return c.JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Helper functions to create response objects (not send them directly)
func ErrorResponse(message string) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
	}
}

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ValidationErrorResponse(errors interface{}) APIResponse {
	return APIResponse{
		Success: false,
		Message: "Validation failed",
		Data:    errors,
	}
}

func SuccessWithPaginationResponse(message string, data interface{}, meta *Meta) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}