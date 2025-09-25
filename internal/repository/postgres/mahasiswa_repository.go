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
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	
	now := time.Now()
	err = sqlDB.QueryRowContext(ctx, query,
		mahasiswa.NIM, mahasiswa.Nama, mahasiswa.Jurusan,
		mahasiswa.Angkatan, mahasiswa.Email, mahasiswa.Password,
		now, now,
	).Scan(&mahasiswa.ID)

	if err != nil {
		return fmt.Errorf("failed to create mahasiswa: %w", err)
	}
	
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
			  FROM mahasiswas WHERE id = $1 AND deleted_at IS NULL`
	
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
			  FROM mahasiswas WHERE nim = $1 AND deleted_at IS NULL`
	
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
			  FROM mahasiswas WHERE email = $1 AND deleted_at IS NULL`
	
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
			  ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
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
	argIndex := 1

	if mahasiswa.NIM != "" {
		setParts = append(setParts, fmt.Sprintf("nim = $%d", argIndex))
		args = append(args, mahasiswa.NIM)
		argIndex++
	}
	if mahasiswa.Nama != "" {
		setParts = append(setParts, fmt.Sprintf("nama = $%d", argIndex))
		args = append(args, mahasiswa.Nama)
		argIndex++
	}
	if mahasiswa.Jurusan != "" {
		setParts = append(setParts, fmt.Sprintf("jurusan = $%d", argIndex))
		args = append(args, mahasiswa.Jurusan)
		argIndex++
	}
	if mahasiswa.Angkatan > 0 {
		setParts = append(setParts, fmt.Sprintf("angkatan = $%d", argIndex))
		args = append(args, mahasiswa.Angkatan)
		argIndex++
	}
	if mahasiswa.Email != "" {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, mahasiswa.Email)
		argIndex++
	}
	if mahasiswa.Password != "" {
		setParts = append(setParts, fmt.Sprintf("password = $%d", argIndex))
		args = append(args, mahasiswa.Password)
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

	query := fmt.Sprintf("UPDATE mahasiswas SET %s WHERE id = $%d AND deleted_at IS NULL", 
		strings.Join(setParts, ", "), argIndex)

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

	query := `UPDATE mahasiswas SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	
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
					 LOWER(nim) LIKE $1 OR LOWER(nama) LIKE $1 OR 
					 LOWER(jurusan) LIKE $1 OR LOWER(email) LIKE $1
				 )`
	var total int64
	err = sqlDB.QueryRowContext(ctx, countSQL, searchQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Get data with pagination
	dataSQL := `SELECT id, nim, nama, jurusan, angkatan, email, password, created_at, updated_at 
				FROM mahasiswas 
				WHERE deleted_at IS NULL AND (
					LOWER(nim) LIKE $1 OR LOWER(nama) LIKE $1 OR 
					LOWER(jurusan) LIKE $1 OR LOWER(email) LIKE $1
				)
				ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := sqlDB.QueryContext(ctx, dataSQL, searchQuery, limit, offset)
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