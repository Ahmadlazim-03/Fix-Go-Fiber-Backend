package repository

import (
	"context"
	"Fix-Go-Fiber-Backend/internal/domain/entity"
)

type AlumniRepository interface {
	Create(ctx context.Context, alumni *entity.Alumni) error
	GetByID(ctx context.Context, id uint) (*entity.Alumni, error)
	GetByMahasiswaID(ctx context.Context, mahasiswaID uint) (*entity.Alumni, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Alumni, int64, error)
	Update(ctx context.Context, id uint, alumni *entity.Alumni) error
	Delete(ctx context.Context, id uint) error
	GetWithMahasiswa(ctx context.Context, id uint) (*entity.Alumni, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*entity.Alumni, int64, error)
}