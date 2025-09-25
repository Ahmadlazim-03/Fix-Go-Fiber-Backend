package service

import (
	"context"
	"Fix-Go-Fiber-Backend/internal/domain/dto"
	"Fix-Go-Fiber-Backend/internal/domain/entity"
)

// PekerjaanAlumniService interface for pekerjaan alumni domain services
type PekerjaanAlumniService interface {
	CreatePekerjaan(ctx context.Context, req *dto.CreatePekerjaanRequest) (*entity.PekerjaanAlumni, error)
	GetPekerjaanByID(ctx context.Context, id uint) (*entity.PekerjaanAlumni, error)
	GetPekerjaanByMahasiswaID(ctx context.Context, mahasiswaID uint) ([]*entity.PekerjaanAlumni, error)
	GetAllPekerjaan(ctx context.Context, search string, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error)
	UpdatePekerjaan(ctx context.Context, id uint, req *dto.UpdatePekerjaanRequest) (*entity.PekerjaanAlumni, error)
	DeletePekerjaan(ctx context.Context, id uint) error
}