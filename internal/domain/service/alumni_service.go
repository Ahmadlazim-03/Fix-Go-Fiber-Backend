package service

import (
	"context"
	"Fix-Go-Fiber-Backend/internal/domain/dto"
	"Fix-Go-Fiber-Backend/internal/domain/entity"
)

// AlumniService interface for alumni domain services
type AlumniService interface {
	CreateAlumni(ctx context.Context, req *dto.CreateAlumniRequest) (*entity.Alumni, error)
	GetAlumniByID(ctx context.Context, id uint) (*entity.Alumni, error)
	GetAlumniByMahasiswaID(ctx context.Context, mahasiswaID uint) (*entity.Alumni, error)
	GetAllAlumni(ctx context.Context, search string, limit, offset int) ([]*entity.Alumni, int64, error)
	UpdateAlumni(ctx context.Context, id uint, req *dto.UpdateAlumniRequest) (*entity.Alumni, error)
	DeleteAlumni(ctx context.Context, id uint) error
}

// PekerjaanAlumniService interface for pekerjaan alumni domain services
type PekerjaanAlumniService interface {
	CreatePekerjaan(ctx context.Context, req *dto.CreatePekerjaanRequest) (*entity.PekerjaanAlumni, error)
	GetPekerjaanByID(ctx context.Context, id uint) (*entity.PekerjaanAlumni, error)
	GetPekerjaanByAlumniID(ctx context.Context, alumniID uint) ([]*entity.PekerjaanAlumni, error)
	GetAllPekerjaan(ctx context.Context, search string, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error)
	UpdatePekerjaan(ctx context.Context, id uint, req *dto.UpdatePekerjaanRequest) (*entity.PekerjaanAlumni, error)
	DeletePekerjaan(ctx context.Context, id uint) error
}