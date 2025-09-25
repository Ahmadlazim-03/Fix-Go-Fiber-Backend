package handler

import (
	"strconv"

	"Fix-Go-Fiber-Backend/internal/domain/dto"
	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
)

type MahasiswaHandler struct {
	mahasiswaUsecase *usecase.MahasiswaUsecase
	validator        *validator.Validate
}

func NewMahasiswaHandler(mahasiswaUsecase *usecase.MahasiswaUsecase, validator *validator.Validate) *MahasiswaHandler {
	return &MahasiswaHandler{
		mahasiswaUsecase: mahasiswaUsecase,
		validator:        validator,
	}
}

func (h *MahasiswaHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateMahasiswaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Validation failed",
			Data:    err.Error(),
		})
	}

	mahasiswa := &entity.Mahasiswa{
		NIM:      req.NIM,
		Nama:     req.Nama,
		Jurusan:  req.Jurusan,
		Angkatan: req.Angkatan,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.mahasiswaUsecase.Create(c.Context(), mahasiswa); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Success: true,
		Message: "Mahasiswa berhasil dibuat",
		Data:    mahasiswa.ToResponse(),
	})
}

func (h *MahasiswaHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid ID",
		})
	}

	mahasiswa, err := h.mahasiswaUsecase.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(dto.APIResponse{
		Success: true,
		Message: "Mahasiswa ditemukan",
		Data:    mahasiswa.ToResponse(),
	})
}

func (h *MahasiswaHandler) GetAll(c *fiber.Ctx) error {
	var query dto.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid query parameters",
		})
	}

	if err := h.validator.Struct(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Validation failed",
			Data:    err.Error(),
		})
	}

	var mahasiswas []*entity.Mahasiswa
	var total int64
	var err error

	if query.Search != "" {
		mahasiswas, total, err = h.mahasiswaUsecase.Search(c.Context(), query.Search, query.Limit, query.GetOffset())
	} else {
		mahasiswas, total, err = h.mahasiswaUsecase.GetAll(c.Context(), query.Limit, query.GetOffset())
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	// Convert to response format
	responses := make([]*entity.MahasiswaResponse, len(mahasiswas))
	for i, m := range mahasiswas {
		responses[i] = m.ToResponse()
	}

	return c.JSON(dto.APIResponse{
		Success: true,
		Message: "Data mahasiswa berhasil diambil",
		Data:    responses,
		Meta:    query.GetMeta(total),
	})
}

func (h *MahasiswaHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid ID",
		})
	}

	var req dto.UpdateMahasiswaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Validation failed",
			Data:    err.Error(),
		})
	}

	mahasiswa := &entity.Mahasiswa{
		NIM:      req.NIM,
		Nama:     req.Nama,
		Jurusan:  req.Jurusan,
		Angkatan: req.Angkatan,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.mahasiswaUsecase.Update(c.Context(), uint(id), mahasiswa); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	// Get updated mahasiswa
	updatedMahasiswa, err := h.mahasiswaUsecase.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to get updated data",
		})
	}

	return c.JSON(dto.APIResponse{
		Success: true,
		Message: "Mahasiswa berhasil diupdate",
		Data:    updatedMahasiswa.ToResponse(),
	})
}

func (h *MahasiswaHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid ID",
		})
	}

	if err := h.mahasiswaUsecase.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(dto.APIResponse{
		Success: true,
		Message: "Mahasiswa berhasil dihapus",
	})
}