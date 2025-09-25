package repository

import (
	"context"
	"Fix-Go-Fiber-Backend/internal/domain/entity"
)

type PekerjaanAlumniRepository interface {
	Create(ctx context.Context, pekerjaan *entity.PekerjaanAlumni) error
	GetByID(ctx context.Context, id uint) (*entity.PekerjaanAlumni, error)
	GetByMahasiswaID(ctx context.Context, mahasiswaID uint) ([]*entity.PekerjaanAlumni, error)
	GetByIDAndMahasiswaID(ctx context.Context, id, mahasiswaID uint) (*entity.PekerjaanAlumni, error)
	GetByMahasiswaIDWithPagination(ctx context.Context, mahasiswaID uint, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error)
	GetAll(ctx context.Context) ([]*entity.PekerjaanAlumni, error)
	GetWithPagination(ctx context.Context, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error)
	Update(ctx context.Context, pekerjaan *entity.PekerjaanAlumni) error
	Delete(ctx context.Context, id uint) error
	GetWithFilters(ctx context.Context, filters map[string]interface{}) ([]*entity.PekerjaanAlumni, error)
}