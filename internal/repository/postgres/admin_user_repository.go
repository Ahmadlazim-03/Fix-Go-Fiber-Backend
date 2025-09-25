package postgres

import (
	"context"

	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/domain/repository"

	"gorm.io/gorm"
)

type adminUserRepository struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) repository.AdminUserRepository {
	return &adminUserRepository{
		db: db,
	}
}

func (r *adminUserRepository) Create(ctx context.Context, admin *entity.AdminUser) error {
	return r.db.WithContext(ctx).Create(admin).Error
}

func (r *adminUserRepository) GetByID(ctx context.Context, id uint) (*entity.AdminUser, error) {
	var admin entity.AdminUser
	err := r.db.WithContext(ctx).First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminUserRepository) GetByUsername(ctx context.Context, username string) (*entity.AdminUser, error) {
	var admin entity.AdminUser
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminUserRepository) GetByEmail(ctx context.Context, email string) (*entity.AdminUser, error) {
	var admin entity.AdminUser
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminUserRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.AdminUser, int64, error) {
	var admins []*entity.AdminUser
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&entity.AdminUser{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get admins with pagination
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&admins).Error

	return admins, total, err
}

func (r *adminUserRepository) Update(ctx context.Context, id uint, admin *entity.AdminUser) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(admin).Error
}

func (r *adminUserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.AdminUser{}, id).Error
}

func (r *adminUserRepository) GetActiveAdmins(ctx context.Context, limit, offset int) ([]*entity.AdminUser, int64, error) {
	var admins []*entity.AdminUser
	var total int64

	// Get total count of active admins
	if err := r.db.WithContext(ctx).Model(&entity.AdminUser{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get active admins with pagination
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Limit(limit).
		Offset(offset).
		Find(&admins).Error

	return admins, total, err
}