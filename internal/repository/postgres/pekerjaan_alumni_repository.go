package postgres

import (
	"context"
	"strings"

	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/domain/repository"

	"gorm.io/gorm"
)

type pekerjaanAlumniRepository struct {
	db *gorm.DB
}

func NewPekerjaanAlumniRepository(db *gorm.DB) repository.PekerjaanAlumniRepository {
	return &pekerjaanAlumniRepository{
		db: db,
	}
}

func (r *pekerjaanAlumniRepository) Create(ctx context.Context, pekerjaan *entity.PekerjaanAlumni) error {
	return r.db.WithContext(ctx).Create(pekerjaan).Error
}

func (r *pekerjaanAlumniRepository) GetByID(ctx context.Context, id uint) (*entity.PekerjaanAlumni, error) {
	var pekerjaan entity.PekerjaanAlumni
	err := r.db.WithContext(ctx).First(&pekerjaan, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &pekerjaan, nil
}

func (r *pekerjaanAlumniRepository) GetByAlumniID(ctx context.Context, alumniID uint) ([]*entity.PekerjaanAlumni, error) {
	var pekerjaan []*entity.PekerjaanAlumni
	err := r.db.WithContext(ctx).
		Where("alumni_id = ?", alumniID).
		Order("created_at DESC").
		Find(&pekerjaan).Error
	return pekerjaan, err
}

func (r *pekerjaanAlumniRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error) {
	var pekerjaan []*entity.PekerjaanAlumni
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&entity.PekerjaanAlumni{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).
		Preload("Alumni").
		Preload("Alumni.Mahasiswa").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&pekerjaan).Error

	return pekerjaan, total, err
}

func (r *pekerjaanAlumniRepository) Update(ctx context.Context, id uint, pekerjaan *entity.PekerjaanAlumni) error {
	updates := make(map[string]interface{})
	
	if pekerjaan.AlumniID != 0 {
		updates["alumni_id"] = pekerjaan.AlumniID
	}
	if pekerjaan.NamaCompany != "" {
		updates["nama_company"] = pekerjaan.NamaCompany
	}
	if pekerjaan.Posisi != "" {
		updates["posisi"] = pekerjaan.Posisi
	}
	if !pekerjaan.TanggalMulai.IsZero() {
		updates["tanggal_mulai"] = pekerjaan.TanggalMulai
	}
	if pekerjaan.TanggalSelesai != nil {
		updates["tanggal_selesai"] = pekerjaan.TanggalSelesai
	}
	if pekerjaan.Status != "" {
		updates["status"] = pekerjaan.Status
	}
	if pekerjaan.Deskripsi != "" {
		updates["deskripsi"] = pekerjaan.Deskripsi
	}

	return r.db.WithContext(ctx).Model(&entity.PekerjaanAlumni{}).Where("id = ?", id).Updates(updates).Error
}

func (r *pekerjaanAlumniRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.PekerjaanAlumni{}, id).Error
}

func (r *pekerjaanAlumniRepository) GetWithAlumni(ctx context.Context, id uint) (*entity.PekerjaanAlumni, error) {
	var pekerjaan entity.PekerjaanAlumni
	err := r.db.WithContext(ctx).
		Preload("Alumni").
		Preload("Alumni.Mahasiswa").
		First(&pekerjaan, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &pekerjaan, nil
}

func (r *pekerjaanAlumniRepository) GetActiveJobs(ctx context.Context, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error) {
	var pekerjaan []*entity.PekerjaanAlumni
	var total int64

	baseQuery := r.db.WithContext(ctx).
		Model(&entity.PekerjaanAlumni{}).
		Where("status = ?", entity.StatusAktif)

	// Count total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := baseQuery.
		Preload("Alumni").
		Preload("Alumni.Mahasiswa").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&pekerjaan).Error

	return pekerjaan, total, err
}

func (r *pekerjaanAlumniRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error) {
	var pekerjaan []*entity.PekerjaanAlumni
	var total int64

	searchQuery := "%" + strings.ToLower(query) + "%"
	
	baseQuery := r.db.WithContext(ctx).
		Model(&entity.PekerjaanAlumni{}).
		Where(
			"LOWER(nama_company) LIKE ? OR LOWER(posisi) LIKE ? OR LOWER(deskripsi) LIKE ?",
			searchQuery, searchQuery, searchQuery,
		)

	// Count total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := baseQuery.
		Preload("Alumni").
		Preload("Alumni.Mahasiswa").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&pekerjaan).Error

	return pekerjaan, total, err
}