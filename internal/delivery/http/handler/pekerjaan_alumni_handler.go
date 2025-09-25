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

type PekerjaanAlumniHandler struct {
	pekerjaanService service.PekerjaanAlumniService
	validator        *validator.Validate
}

func NewPekerjaanAlumniHandler(pekerjaanService service.PekerjaanAlumniService, validator *validator.Validate) *PekerjaanAlumniHandler {
	return &PekerjaanAlumniHandler{
		pekerjaanService: pekerjaanService,
		validator:        validator,
	}
}

// CreatePekerjaan - Alumni for self, Admin for any
func (h *PekerjaanAlumniHandler) CreatePekerjaan(c *fiber.Ctx) error {
	var req dto.CreatePekerjaanRequest
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
	
	// If user is alumni, they can only create pekerjaan for themselves
	if claims.Role == "alumni" {
		if req.AlumniID != claims.UserID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Access denied: You can only create pekerjaan for yourself",
			})
		}
	}

	pekerjaan, err := h.pekerjaanService.CreatePekerjaan(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create pekerjaan",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan created successfully",
		"data":    pekerjaan.ToResponse(),
	})
}

// GetAllPekerjaan - Admin only
func (h *PekerjaanAlumniHandler) GetAllPekerjaan(c *fiber.Ctx) error {
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

	pekerjaan, total, err := h.pekerjaanService.GetAllPekerjaan(c.Context(), query, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get pekerjaan",
			"error":   err.Error(),
		})
	}

	response := make([]*entity.PekerjaanAlumniResponse, len(pekerjaan))
	for i, p := range pekerjaan {
		response[i] = p.ToResponse()
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan retrieved successfully",
		"data": fiber.Map{
			"pekerjaan":   response,
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetPekerjaanByAlumniID - Alumni for self, Admin for any
func (h *PekerjaanAlumniHandler) GetPekerjaanByAlumniID(c *fiber.Ctx) error {
	alumniIDStr := c.Params("alumni_id")
	alumniID, err := strconv.ParseUint(alumniIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid alumni ID",
		})
	}

	// Get user claims from JWT
	claims := c.Locals("user").(*jwt.Claims)
	
	// If user is alumni, check if they can only access their own pekerjaan
	if claims.Role == "alumni" {
		if claims.UserID != uint(alumniID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Access denied: You can only view your own pekerjaan",
			})
		}
	}

	pekerjaan, err := h.pekerjaanService.GetPekerjaanByAlumniID(c.Context(), uint(alumniID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get pekerjaan",
			"error":   err.Error(),
		})
	}

	response := make([]*entity.PekerjaanAlumniResponse, len(pekerjaan))
	for i, p := range pekerjaan {
		response[i] = p.ToResponse()
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan retrieved successfully",
		"data":    response,
	})
}

// GetPekerjaanByID - Alumni for own, Admin for any
func (h *PekerjaanAlumniHandler) GetPekerjaanByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid pekerjaan ID",
		})
	}

	pekerjaan, err := h.pekerjaanService.GetPekerjaanByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Pekerjaan not found",
			"error":   err.Error(),
		})
	}

	// Get user claims from JWT
	claims := c.Locals("user").(*jwt.Claims)
	
	// If user is alumni, check if they can only access their own pekerjaan
	if claims.Role == "alumni" {
		if claims.UserID != pekerjaan.AlumniID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Access denied: You can only view your own pekerjaan",
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan retrieved successfully",
		"data":    pekerjaan.ToResponse(),
	})
}

// UpdatePekerjaan - Alumni for own, Admin for any
func (h *PekerjaanAlumniHandler) UpdatePekerjaan(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid pekerjaan ID",
		})
	}

	var req dto.UpdatePekerjaanRequest
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

	// Check if pekerjaan exists and get owner info
	existingPekerjaan, err := h.pekerjaanService.GetPekerjaanByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Pekerjaan not found",
			"error":   err.Error(),
		})
	}

	// Get user claims from JWT
	claims := c.Locals("user").(*jwt.Claims)
	
	// If user is alumni, check if they can only update their own pekerjaan
	if claims.Role == "alumni" {
		if claims.UserID != existingPekerjaan.AlumniID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Access denied: You can only update your own pekerjaan",
			})
		}
	}

	pekerjaan, err := h.pekerjaanService.UpdatePekerjaan(c.Context(), uint(id), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update pekerjaan",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan updated successfully",
		"data":    pekerjaan.ToResponse(),
	})
}

// DeletePekerjaan - Alumni for own, Admin for any (soft delete)
func (h *PekerjaanAlumniHandler) DeletePekerjaan(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid pekerjaan ID",
		})
	}

	// Check if pekerjaan exists and get owner info
	existingPekerjaan, err := h.pekerjaanService.GetPekerjaanByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Pekerjaan not found",
			"error":   err.Error(),
		})
	}

	// Get user claims from JWT
	claims := c.Locals("user").(*jwt.Claims)
	
	// If user is alumni, check if they can only delete their own pekerjaan
	if claims.Role == "alumni" {
		if claims.UserID != existingPekerjaan.AlumniID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Access denied: You can only delete your own pekerjaan",
			})
		}
	}

	err = h.pekerjaanService.DeletePekerjaan(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete pekerjaan",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan soft deleted successfully",
	})
}