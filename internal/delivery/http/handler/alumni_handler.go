package handler

import (
	"Fix-Go-Fiber-Backend/internal/domain/dto"
	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/domain/service"
	"Fix-Go-Fiber-Backend/pkg/jwt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AlumniHandler struct {
	alumniService service.AlumniService
	validator     *validator.Validate
}

func NewAlumniHandler(alumniService service.AlumniService, validator *validator.Validate) *AlumniHandler {
	return &AlumniHandler{
		alumniService: alumniService,
		validator:     validator,
	}
}

// CreateAlumni - Admin only
func (h *AlumniHandler) CreateAlumni(c *fiber.Ctx) error {
	var req dto.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	alumni, err := h.alumniService.CreateAlumni(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create alumni",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Alumni created successfully",
		"data":    alumni.ToResponse(),
	})
}

// GetAllAlumni - Admin only
func (h *AlumniHandler) GetAllAlumni(c *fiber.Ctx) error {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	offset := (page - 1) * limit
	query := c.Query("search")

	alumni, total, err := h.alumniService.GetAllAlumni(c.Context(), query, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get alumni",
			"error":   err.Error(),
		})
	}

	response := make([]*entity.AlumniResponse, len(alumni))
	for i, a := range alumni {
		response[i] = a.ToResponse()
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni retrieved successfully",
		"data": fiber.Map{
			"alumni":      response,
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetAlumniByID - Admin and Self (Alumni)
func (h *AlumniHandler) GetAlumniByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid alumni ID",
		})
	}

	// Get user claims from JWT
	claims := c.Locals("user").(*jwt.Claims)
	
	// If user is alumni, check if they can only access their own data
	if claims.Role == "alumni" {
		if claims.UserID != uint(id) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Access denied: You can only view your own alumni data",
			})
		}
	}

	alumni, err := h.alumniService.GetAlumniByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Alumni not found",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni retrieved successfully",
		"data":    alumni.ToResponse(),
	})
}

// UpdateAlumni - Admin and Self (Alumni)
func (h *AlumniHandler) UpdateAlumni(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid alumni ID",
		})
	}

	var req dto.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	// Get user claims from JWT
	claims := c.Locals("user").(*jwt.Claims)
	
	// If user is alumni, check if they can only update their own data
	if claims.Role == "alumni" {
		if claims.UserID != uint(id) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Access denied: You can only update your own alumni data",
			})
		}
	}

	alumni, err := h.alumniService.UpdateAlumni(c.Context(), uint(id), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update alumni",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni updated successfully",
		"data":    alumni.ToResponse(),
	})
}

// DeleteAlumni - Admin only (soft delete)
func (h *AlumniHandler) DeleteAlumni(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid alumni ID",
		})
	}

	// Only admin can soft delete alumni
	err = h.alumniService.DeleteAlumni(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete alumni",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni soft deleted successfully",
	})
}