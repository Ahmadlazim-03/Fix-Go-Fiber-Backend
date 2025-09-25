package postgres

import (
	"context"
	"strings"

	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/domain/repository"

	"gorm.io/gorm"
)

type mahasiswaRepository struct {
	db *gorm.DB
}

func NewMahasiswaRepository(db *gorm.DB) repository.MahasiswaRepository {
	return &mahasiswaRepository{
		db: db,
	}
}

func (r *mahasiswaRepository) Create(ctx context.Context, mahasiswa *entity.Mahasiswa) error {
	return r.db.WithContext(ctx).Create(mahasiswa).Error
}

func (r *mahasiswaRepository) GetByID(ctx context.Context, id uint) (*entity.Mahasiswa, error) {
	var mahasiswa entity.Mahasiswa
	err := r.db.WithContext(ctx).First(&mahasiswa, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &mahasiswa, nil
}

func (r *mahasiswaRepository) GetByNIM(ctx context.Context, nim string) (*entity.Mahasiswa, error) {
	var mahasiswa entity.Mahasiswa
	err := r.db.WithContext(ctx).Where("nim = ?", nim).First(&mahasiswa).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &mahasiswa, nil
}

func (r *mahasiswaRepository) GetByEmail(ctx context.Context, email string) (*entity.Mahasiswa, error) {
	var mahasiswa entity.Mahasiswa
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&mahasiswa).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &mahasiswa, nil
}

// GetByEmail without context for auth service compatibility
func (r *mahasiswaRepository) GetByEmailSimple(email string) (*entity.Mahasiswa, error) {
	var mahasiswa entity.Mahasiswa
	err := r.db.Where("email = ?", email).First(&mahasiswa).Error
	if err != nil {
		return nil, err
	}
	return &mahasiswa, nil
}

func (r *mahasiswaRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Mahasiswa, int64, error) {
	var mahasiswas []*entity.Mahasiswa
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&entity.Mahasiswa{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&mahasiswas).Error

	return mahasiswas, total, err
}

func (r *mahasiswaRepository) Update(ctx context.Context, id uint, mahasiswa *entity.Mahasiswa) error {
	updates := make(map[string]interface{})
	
	if mahasiswa.NIM != "" {
		updates["nim"] = mahasiswa.NIM
	}
	if mahasiswa.Nama != "" {
		updates["nama"] = mahasiswa.Nama
	}
	if mahasiswa.Jurusan != "" {
		updates["jurusan"] = mahasiswa.Jurusan
	}
	if mahasiswa.Angkatan > 0 {
		updates["angkatan"] = mahasiswa.Angkatan
	}
	if mahasiswa.Email != "" {
		updates["email"] = mahasiswa.Email
	}
	if mahasiswa.Password != "" {
		updates["password"] = mahasiswa.Password
	}

	return r.db.WithContext(ctx).Model(&entity.Mahasiswa{}).Where("id = ?", id).Updates(updates).Error
}

func (r *mahasiswaRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Mahasiswa{}, id).Error
}

func (r *mahasiswaRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.Mahasiswa, int64, error) {
	var mahasiswas []*entity.Mahasiswa
	var total int64

	searchQuery := "%" + strings.ToLower(query) + "%"
	
	baseQuery := r.db.WithContext(ctx).Model(&entity.Mahasiswa{}).Where(
		"LOWER(nim) LIKE ? OR LOWER(nama) LIKE ? OR LOWER(jurusan) LIKE ? OR LOWER(email) LIKE ?",
		searchQuery, searchQuery, searchQuery, searchQuery,
	)

	// Count total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := baseQuery.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&mahasiswas).Error

	return mahasiswas, total, err
}