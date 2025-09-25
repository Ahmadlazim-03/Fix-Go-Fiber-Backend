package repository

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

type mahasiswaRepository struct {
	db *gorm.DB
}

func NewMahasiswaRepository(db *gorm.DB) repository.MahasiswaRepository {
	return &mahasiswaRepository{
		db: db,
	}
}

func (r *mahasiswaRepository) Create(ctx context.Context, mahasiswa *entity.Mahasiswa) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	query := `INSERT INTO mahasiswas (nim, nama, jurusan, angkatan, email, password, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	result, err := sqlDB.ExecContext(ctx, query,
		mahasiswa.NIM, mahasiswa.Nama, mahasiswa.Jurusan,
		mahasiswa.Angkatan, mahasiswa.Email, mahasiswa.Password,
		now, now,
	)

	if err != nil {
		return fmt.Errorf("failed to create mahasiswa: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	
	mahasiswa.ID = uint(id)
	mahasiswa.CreatedAt = now
	mahasiswa.UpdatedAt = now
	return nil
}

func (r *mahasiswaRepository) GetByID(ctx context.Context, id uint) (*entity.Mahasiswa, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, nim, nama, jurusan, angkatan, email, password, created_at, updated_at 
			  FROM mahasiswas WHERE id = ? AND deleted_at IS NULL`
	
	var mahasiswa entity.Mahasiswa
	err = sqlDB.QueryRowContext(ctx, query, id).Scan(
		&mahasiswa.ID, &mahasiswa.NIM, &mahasiswa.Nama,
		&mahasiswa.Jurusan, &mahasiswa.Angkatan, &mahasiswa.Email,
		&mahasiswa.Password, &mahasiswa.CreatedAt, &mahasiswa.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get mahasiswa by ID: %w", err)
	}

	return &mahasiswa, nil
}

func (r *mahasiswaRepository) GetByNIM(ctx context.Context, nim string) (*entity.Mahasiswa, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, nim, nama, jurusan, angkatan, email, password, created_at, updated_at 
			  FROM mahasiswas WHERE nim = ? AND deleted_at IS NULL`
	
	var mahasiswa entity.Mahasiswa
	err = sqlDB.QueryRowContext(ctx, query, nim).Scan(
		&mahasiswa.ID, &mahasiswa.NIM, &mahasiswa.Nama,
		&mahasiswa.Jurusan, &mahasiswa.Angkatan, &mahasiswa.Email,
		&mahasiswa.Password, &mahasiswa.CreatedAt, &mahasiswa.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get mahasiswa by NIM: %w", err)
	}

	return &mahasiswa, nil
}

func (r *mahasiswaRepository) GetByEmail(ctx context.Context, email string) (*entity.Mahasiswa, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, nim, nama, jurusan, angkatan, email, password, created_at, updated_at 
			  FROM mahasiswas WHERE email = ? AND deleted_at IS NULL`
	
	var mahasiswa entity.Mahasiswa
	err = sqlDB.QueryRowContext(ctx, query, email).Scan(
		&mahasiswa.ID, &mahasiswa.NIM, &mahasiswa.Nama,
		&mahasiswa.Jurusan, &mahasiswa.Angkatan, &mahasiswa.Email,
		&mahasiswa.Password, &mahasiswa.CreatedAt, &mahasiswa.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get mahasiswa by email: %w", err)
	}

	return &mahasiswa, nil
}

// GetByEmail without context for auth service compatibility
func (r *mahasiswaRepository) GetByEmailSimple(email string) (*entity.Mahasiswa, error) {
	return r.GetByEmail(context.Background(), email)
}

func (r *mahasiswaRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Mahasiswa, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM mahasiswas WHERE deleted_at IS NULL`
	var total int64
	err = sqlDB.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count mahasiswa: %w", err)
	}

	// Get data with pagination
	query := `SELECT id, nim, nama, jurusan, angkatan, email, password, created_at, updated_at 
			  FROM mahasiswas WHERE deleted_at IS NULL 
			  ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	rows, err := sqlDB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get mahasiswa list: %w", err)
	}
	defer rows.Close()

	var mahasiswas []*entity.Mahasiswa
	for rows.Next() {
		var mahasiswa entity.Mahasiswa
		err = rows.Scan(
			&mahasiswa.ID, &mahasiswa.NIM, &mahasiswa.Nama,
			&mahasiswa.Jurusan, &mahasiswa.Angkatan, &mahasiswa.Email,
			&mahasiswa.Password, &mahasiswa.CreatedAt, &mahasiswa.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan mahasiswa: %w", err)
		}
		mahasiswas = append(mahasiswas, &mahasiswa)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating mahasiswa rows: %w", err)
	}

	return mahasiswas, total, nil
}

func (r *mahasiswaRepository) Update(ctx context.Context, id uint, mahasiswa *entity.Mahasiswa) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	setParts := []string{}
	args := []interface{}{}

	if mahasiswa.NIM != "" {
		setParts = append(setParts, "nim = ?")
		args = append(args, mahasiswa.NIM)
	}
	if mahasiswa.Nama != "" {
		setParts = append(setParts, "nama = ?")
		args = append(args, mahasiswa.Nama)
	}
	if mahasiswa.Jurusan != "" {
		setParts = append(setParts, "jurusan = ?")
		args = append(args, mahasiswa.Jurusan)
	}
	if mahasiswa.Angkatan > 0 {
		setParts = append(setParts, "angkatan = ?")
		args = append(args, mahasiswa.Angkatan)
	}
	if mahasiswa.Email != "" {
		setParts = append(setParts, "email = ?")
		args = append(args, mahasiswa.Email)
	}
	if mahasiswa.Password != "" {
		setParts = append(setParts, "password = ?")
		args = append(args, mahasiswa.Password)
	}

	if len(setParts) == 0 {
		return errors.New("no fields to update")
	}

	// Add updated_at
	setParts = append(setParts, "updated_at = ?")
	args = append(args, time.Now())

	// Add WHERE condition
	args = append(args, id)

	query := fmt.Sprintf("UPDATE mahasiswas SET %s WHERE id = ? AND deleted_at IS NULL", 
		strings.Join(setParts, ", "))

	result, err := sqlDB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update mahasiswa: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("mahasiswa not found or already deleted")
	}

	return nil
}

func (r *mahasiswaRepository) Delete(ctx context.Context, id uint) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	query := `UPDATE mahasiswas SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL`
	
	result, err := sqlDB.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete mahasiswa: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("mahasiswa not found or already deleted")
	}

	return nil
}

func (r *mahasiswaRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.Mahasiswa, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	searchQuery := "%" + strings.ToLower(query) + "%"

	// Count total
	countSQL := `SELECT COUNT(*) FROM mahasiswas 
				 WHERE deleted_at IS NULL AND (
					 LOWER(nim) LIKE ? OR LOWER(nama) LIKE ? OR 
					 LOWER(jurusan) LIKE ? OR LOWER(email) LIKE ?
				 )`
	var total int64
	err = sqlDB.QueryRowContext(ctx, countSQL, searchQuery, searchQuery, searchQuery, searchQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Get data with pagination
	dataSQL := `SELECT id, nim, nama, jurusan, angkatan, email, password, created_at, updated_at 
				FROM mahasiswas 
				WHERE deleted_at IS NULL AND (
					LOWER(nim) LIKE ? OR LOWER(nama) LIKE ? OR 
					LOWER(jurusan) LIKE ? OR LOWER(email) LIKE ?
				)
				ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	rows, err := sqlDB.QueryContext(ctx, dataSQL, searchQuery, searchQuery, searchQuery, searchQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search mahasiswa: %w", err)
	}
	defer rows.Close()

	var mahasiswas []*entity.Mahasiswa
	for rows.Next() {
		var mahasiswa entity.Mahasiswa
		err = rows.Scan(
			&mahasiswa.ID, &mahasiswa.NIM, &mahasiswa.Nama,
			&mahasiswa.Jurusan, &mahasiswa.Angkatan, &mahasiswa.Email,
			&mahasiswa.Password, &mahasiswa.CreatedAt, &mahasiswa.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan mahasiswa: %w", err)
		}
		mahasiswas = append(mahasiswas, &mahasiswa)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating search results: %w", err)
	}

	return mahasiswas, total, nil
}