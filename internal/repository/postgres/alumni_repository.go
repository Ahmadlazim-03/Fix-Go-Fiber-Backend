package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	query := `INSERT INTO alumni (mahasiswa_id, tahun_lulus, no_telepon, alamat, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	
	now := time.Now()
	err = sqlDB.QueryRowContext(ctx, query,
		alumni.MahasiswaID, alumni.TahunLulus, alumni.NoTelepon,
		alumni.Alamat, now, now,
	).Scan(&alumni.ID)

	if err != nil {
		return fmt.Errorf("failed to create alumni: %w", err)
	}
	
	alumni.CreatedAt = now
	alumni.UpdatedAt = now
	return nil
}

func (r *alumniRepository) GetByID(ctx context.Context, id uint) (*entity.Alumni, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, mahasiswa_id, tahun_lulus, no_telepon, alamat, created_at, updated_at 
			  FROM alumni WHERE id = $1 AND deleted_at IS NULL`
	
	var alumni entity.Alumni
	err = sqlDB.QueryRowContext(ctx, query, id).Scan(
		&alumni.ID, &alumni.MahasiswaID, &alumni.TahunLulus,
		&alumni.NoTelepon, &alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get alumni by ID: %w", err)
	}

	return &alumni, nil
}

func (r *alumniRepository) GetByMahasiswaID(ctx context.Context, mahasiswaID uint) (*entity.Alumni, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, mahasiswa_id, tahun_lulus, no_telepon, alamat, created_at, updated_at 
			  FROM alumni WHERE mahasiswa_id = $1 AND deleted_at IS NULL`
	
	var alumni entity.Alumni
	err = sqlDB.QueryRowContext(ctx, query, mahasiswaID).Scan(
		&alumni.ID, &alumni.MahasiswaID, &alumni.TahunLulus,
		&alumni.NoTelepon, &alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get alumni by mahasiswa ID: %w", err)
	}

	return &alumni, nil
}

func (r *alumniRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Alumni, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM alumni WHERE deleted_at IS NULL`
	var total int64
	err = sqlDB.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count alumni: %w", err)
	}

	// Get data with mahasiswa info (JOIN)
	query := `SELECT a.id, a.mahasiswa_id, a.tahun_lulus, a.no_telepon, a.alamat, a.created_at, a.updated_at,
					 m.id, m.nim, m.nama, m.jurusan, m.angkatan, m.email, m.password, m.created_at, m.updated_at
			  FROM alumni a
			  JOIN mahasiswas m ON a.mahasiswa_id = m.id
			  WHERE a.deleted_at IS NULL AND m.deleted_at IS NULL
			  ORDER BY a.created_at DESC LIMIT $1 OFFSET $2`
	
	rows, err := sqlDB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get alumni list: %w", err)
	}
	defer rows.Close()

	var alumni []*entity.Alumni
	for rows.Next() {
		var alum entity.Alumni
		var mahasiswa entity.Mahasiswa
		
		err = rows.Scan(
			&alum.ID, &alum.MahasiswaID, &alum.TahunLulus, &alum.NoTelepon, &alum.Alamat, &alum.CreatedAt, &alum.UpdatedAt,
			&mahasiswa.ID, &mahasiswa.NIM, &mahasiswa.Nama, &mahasiswa.Jurusan, &mahasiswa.Angkatan, 
			&mahasiswa.Email, &mahasiswa.Password, &mahasiswa.CreatedAt, &mahasiswa.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan alumni: %w", err)
		}
		
		alum.Mahasiswa = mahasiswa
		alumni = append(alumni, &alum)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating alumni rows: %w", err)
	}

	return alumni, total, nil
}

func (r *alumniRepository) Update(ctx context.Context, id uint, alumni *entity.Alumni) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if alumni.MahasiswaID != 0 {
		setParts = append(setParts, fmt.Sprintf("mahasiswa_id = $%d", argIndex))
		args = append(args, alumni.MahasiswaID)
		argIndex++
	}
	if alumni.TahunLulus > 0 {
		setParts = append(setParts, fmt.Sprintf("tahun_lulus = $%d", argIndex))
		args = append(args, alumni.TahunLulus)
		argIndex++
	}
	if alumni.NoTelepon != "" {
		setParts = append(setParts, fmt.Sprintf("no_telepon = $%d", argIndex))
		args = append(args, alumni.NoTelepon)
		argIndex++
	}
	if alumni.Alamat != "" {
		setParts = append(setParts, fmt.Sprintf("alamat = $%d", argIndex))
		args = append(args, alumni.Alamat)
		argIndex++
	}

	if len(setParts) == 0 {
		return errors.New("no fields to update")
	}

	// Add updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE condition
	args = append(args, id)

	query := fmt.Sprintf("UPDATE alumni SET %s WHERE id = $%d AND deleted_at IS NULL", 
		strings.Join(setParts, ", "), argIndex)

	result, err := sqlDB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update alumni: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("alumni not found or already deleted")
	}

	return nil
}

func (r *alumniRepository) Delete(ctx context.Context, id uint) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	query := `UPDATE alumni SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	
	result, err := sqlDB.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete alumni: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("alumni not found or already deleted")
	}

	return nil
}

func (r *alumniRepository) GetWithMahasiswa(ctx context.Context, id uint) (*entity.Alumni, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT a.id, a.mahasiswa_id, a.tahun_lulus, a.no_telepon, a.alamat, a.created_at, a.updated_at,
					 m.id, m.nim, m.nama, m.jurusan, m.angkatan, m.email, m.password, m.created_at, m.updated_at
			  FROM alumni a
			  JOIN mahasiswas m ON a.mahasiswa_id = m.id
			  WHERE a.id = $1 AND a.deleted_at IS NULL AND m.deleted_at IS NULL`
	
	var alum entity.Alumni
	var mahasiswa entity.Mahasiswa
	
	err = sqlDB.QueryRowContext(ctx, query, id).Scan(
		&alum.ID, &alum.MahasiswaID, &alum.TahunLulus, &alum.NoTelepon, &alum.Alamat, &alum.CreatedAt, &alum.UpdatedAt,
		&mahasiswa.ID, &mahasiswa.NIM, &mahasiswa.Nama, &mahasiswa.Jurusan, &mahasiswa.Angkatan, 
		&mahasiswa.Email, &mahasiswa.Password, &mahasiswa.CreatedAt, &mahasiswa.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get alumni with mahasiswa: %w", err)
	}

	alum.Mahasiswa = mahasiswa
	return &alum, nil
}

func (r *alumniRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.Alumni, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	searchQuery := "%" + strings.ToLower(query) + "%"

	// Count total
	countSQL := `SELECT COUNT(*) FROM alumni a
				 JOIN mahasiswas m ON a.mahasiswa_id = m.id
				 WHERE a.deleted_at IS NULL AND m.deleted_at IS NULL AND (
					 LOWER(m.nama) LIKE $1 OR LOWER(m.nim) LIKE $1 OR LOWER(a.alamat) LIKE $1
				 )`
	var total int64
	err = sqlDB.QueryRowContext(ctx, countSQL, searchQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Get data with mahasiswa info
	dataSQL := `SELECT a.id, a.mahasiswa_id, a.tahun_lulus, a.no_telepon, a.alamat, a.created_at, a.updated_at,
					   m.id, m.nim, m.nama, m.jurusan, m.angkatan, m.email, m.password, m.created_at, m.updated_at
				FROM alumni a
				JOIN mahasiswas m ON a.mahasiswa_id = m.id
				WHERE a.deleted_at IS NULL AND m.deleted_at IS NULL AND (
					LOWER(m.nama) LIKE $1 OR LOWER(m.nim) LIKE $1 OR LOWER(a.alamat) LIKE $1
				)
				ORDER BY a.created_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := sqlDB.QueryContext(ctx, dataSQL, searchQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search alumni: %w", err)
	}
	defer rows.Close()

	var alumni []*entity.Alumni
	for rows.Next() {
		var alum entity.Alumni
		var mahasiswa entity.Mahasiswa
		
		err = rows.Scan(
			&alum.ID, &alum.MahasiswaID, &alum.TahunLulus, &alum.NoTelepon, &alum.Alamat, &alum.CreatedAt, &alum.UpdatedAt,
			&mahasiswa.ID, &mahasiswa.NIM, &mahasiswa.Nama, &mahasiswa.Jurusan, &mahasiswa.Angkatan, 
			&mahasiswa.Email, &mahasiswa.Password, &mahasiswa.CreatedAt, &mahasiswa.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan alumni: %w", err)
		}
		
		alum.Mahasiswa = mahasiswa
		alumni = append(alumni, &alum)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating search results: %w", err)
	}

	return alumni, total, nil
}