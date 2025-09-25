package repository

import (
	"context"
	"Fix-Go-Fiber-Backend/internal/domain/entity"
)

type AdminUserRepository interface {
	Create(ctx context.Context, admin *entity.AdminUser) error
	GetByID(ctx context.Context, id uint) (*entity.AdminUser, error)
	GetByUsername(ctx context.Context, username string) (*entity.AdminUser, error)
	GetByEmail(ctx context.Context, email string) (*entity.AdminUser, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entity.AdminUser, int64, error)
	Update(ctx context.Context, id uint, admin *entity.AdminUser) error
	Delete(ctx context.Context, id uint) error
	GetActiveAdmins(ctx context.Context, limit, offset int) ([]*entity.AdminUser, int64, error)
}