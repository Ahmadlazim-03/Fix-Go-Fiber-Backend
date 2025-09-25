package handler

import (
	"Fix-Go-Fiber-Backend/internal/domain/dto"
	"Fix-Go-Fiber-Backend/internal/domain/service"
	"Fix-Go-Fiber-Backend/pkg/utils"
	"Fix-Go-Fiber-Backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
	validator   *validator.CustomValidator
}

func NewAuthHandler(authService service.AuthService, validator *validator.CustomValidator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
	}
}

// LoginMahasiswa handles mahasiswa login
func (h *AuthHandler) LoginMahasiswa(c *fiber.Ctx) error {
	var req dto.MahasiswaLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidationErrorResponse(err))
	}

	response, err := h.authService.LoginMahasiswa(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse(err.Error()))
	}

	return c.JSON(utils.SuccessResponse("Login successful", response))
}

// LoginAlumni handles alumni login
func (h *AuthHandler) LoginAlumni(c *fiber.Ctx) error {
	var req dto.AlumniLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidationErrorResponse(err))
	}

	response, err := h.authService.LoginAlumni(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse(err.Error()))
	}

	return c.JSON(utils.SuccessResponse("Login successful", response))
}

// RegisterMahasiswa handles mahasiswa registration
func (h *AuthHandler) RegisterMahasiswa(c *fiber.Ctx) error {
	var req dto.RegisterMahasiswaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidationErrorResponse(err))
	}

	response, err := h.authService.RegisterMahasiswa(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse("Registration successful", response))
}

// GraduateMahasiswa handles marking mahasiswa as graduated (alumni)
func (h *AuthHandler) GraduateMahasiswa(c *fiber.Ctx) error {
	var req dto.GraduateMahasiswaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidationErrorResponse(err))
	}

	response, err := h.authService.GraduateMahasiswa(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Graduation successful", response))
}

// LoginAdmin handles admin login
func (h *AuthHandler) LoginAdmin(c *fiber.Ctx) error {
	var req dto.AdminLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidationErrorResponse(err))
	}

	response, err := h.authService.LoginAdmin(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse(err.Error()))
	}

	return c.JSON(utils.SuccessResponse("Login successful", response))
}

// GetProfile returns current user profile based on token
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(uint)
	email := c.Locals("email").(string)

	profile := map[string]interface{}{
		"id":    userID,
		"email": email,
		"role":  role,
	}

	if username := c.Locals("username"); username != nil {
		profile["username"] = username
	}

	return c.JSON(utils.SuccessResponse("Profile retrieved successfully", profile))
}