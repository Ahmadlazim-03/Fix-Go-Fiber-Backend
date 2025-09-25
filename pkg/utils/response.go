package utils

import (
	"Fix-Go-Fiber-Backend/internal/delivery/http/dto"

	"github.com/gofiber/fiber/v2"
)

func SendSuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SendErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(dto.APIResponse{
		Success: false,
		Message: message,
	})
}

func SendValidationErrorResponse(c *fiber.Ctx, message string, errors interface{}) error {
	return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
		Success: false,
		Message: message,
		Data:    errors,
	})
}

func SendPaginatedResponse(c *fiber.Ctx, message string, data interface{}, meta *dto.Meta) error {
	return c.JSON(dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Helper functions to create response objects (not send them directly)
func ErrorResponse(message string) dto.APIResponse {
	return dto.APIResponse{
		Success: false,
		Message: message,
	}
}

func SuccessResponse(message string, data interface{}) dto.APIResponse {
	return dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ValidationErrorResponse(errors interface{}) dto.APIResponse {
	return dto.APIResponse{
		Success: false,
		Message: "Validation failed",
		Data:    errors,
	}
}

func SuccessWithPaginationResponse(message string, data interface{}, meta *dto.Meta) dto.APIResponse {
	return dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}