package repository

import (
	"context"
	"Fix-Go-Fiber-Backend/internal/domain/entity"
)

type PekerjaanAlumniRepository interface {
	Create(ctx context.Context, pekerjaan *entity.PekerjaanAlumni) error
	GetByID(ctx context.Context, id uint) (*entity.PekerjaanAlumni, error)
	GetByAlumniID(ctx context.Context, alumniID uint) ([]*entity.PekerjaanAlumni, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error)
	Update(ctx context.Context, id uint, pekerjaan *entity.PekerjaanAlumni) error
	Delete(ctx context.Context, id uint) error
	GetWithAlumni(ctx context.Context, id uint) (*entity.PekerjaanAlumni, error)
	GetActiveJobs(ctx context.Context, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error)
}