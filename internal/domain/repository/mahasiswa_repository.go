package repository

import (
	"context"
	"Fix-Go-Fiber-Backend/internal/domain/entity"
)

type MahasiswaRepository interface {
	Create(ctx context.Context, mahasiswa *entity.Mahasiswa) error
	GetByID(ctx context.Context, id uint) (*entity.Mahasiswa, error)
	GetByNIM(ctx context.Context, nim string) (*entity.Mahasiswa, error)
	GetByEmail(ctx context.Context, email string) (*entity.Mahasiswa, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Mahasiswa, int64, error)
	Update(ctx context.Context, id uint, mahasiswa *entity.Mahasiswa) error
	Delete(ctx context.Context, id uint) error
	Search(ctx context.Context, query string, limit, offset int) ([]*entity.Mahasiswa, int64, error)
}