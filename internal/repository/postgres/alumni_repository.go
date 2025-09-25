package postgres

import (
	"context"
	"strings"

	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/domain/repository"

	"gorm.io/gorm"
)

type alumniRepository struct {
	db *gorm.DB
}

func NewAlumniRepository(db *gorm.DB) repository.AlumniRepository {
	return &alumniRepository{
		db: db,
	}
}

func (r *alumniRepository) Create(ctx context.Context, alumni *entity.Alumni) error {
	return r.db.WithContext(ctx).Create(alumni).Error
}

func (r *alumniRepository) GetByID(ctx context.Context, id uint) (*entity.Alumni, error) {
	var alumni entity.Alumni
	err := r.db.WithContext(ctx).First(&alumni, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &alumni, nil
}

func (r *alumniRepository) GetByMahasiswaID(ctx context.Context, mahasiswaID uint) (*entity.Alumni, error) {
	var alumni entity.Alumni
	err := r.db.WithContext(ctx).Where("mahasiswa_id = ?", mahasiswaID).First(&alumni).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &alumni, nil
}

func (r *alumniRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Alumni, int64, error) {
	var alumni []*entity.Alumni
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&entity.Alumni{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).
		Preload("Mahasiswa").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&alumni).Error

	return alumni, total, err
}

func (r *alumniRepository) Update(ctx context.Context, id uint, alumni *entity.Alumni) error {
	updates := make(map[string]interface{})
	
	if alumni.MahasiswaID != 0 {
		updates["mahasiswa_id"] = alumni.MahasiswaID
	}
	if alumni.TahunLulus > 0 {
		updates["tahun_lulus"] = alumni.TahunLulus
	}
	if alumni.NoTelepon != "" {
		updates["no_telepon"] = alumni.NoTelepon
	}
	if alumni.Alamat != "" {
		updates["alamat"] = alumni.Alamat
	}

	return r.db.WithContext(ctx).Model(&entity.Alumni{}).Where("id = ?", id).Updates(updates).Error
}

func (r *alumniRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Alumni{}, id).Error
}

func (r *alumniRepository) GetWithMahasiswa(ctx context.Context, id uint) (*entity.Alumni, error) {
	var alumni entity.Alumni
	err := r.db.WithContext(ctx).
		Preload("Mahasiswa").
		First(&alumni, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &alumni, nil
}

func (r *alumniRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.Alumni, int64, error) {
	var alumni []*entity.Alumni
	var total int64

	searchQuery := "%" + strings.ToLower(query) + "%"
	
	baseQuery := r.db.WithContext(ctx).
		Model(&entity.Alumni{}).
		Joins("JOIN mahasiswa ON alumni.mahasiswa_id = mahasiswa.id").
		Where(
			"LOWER(mahasiswa.nama) LIKE ? OR LOWER(mahasiswa.nim) LIKE ? OR LOWER(alumni.alamat) LIKE ?",
			searchQuery, searchQuery, searchQuery,
		)

	// Count total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := baseQuery.
		Preload("Mahasiswa").
		Order("alumni.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&alumni).Error

	return alumni, total, err
}