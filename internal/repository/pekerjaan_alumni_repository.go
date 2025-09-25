package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/domain/repository"

	"gorm.io/gorm"
)

type pekerjaanAlumniRepository struct {
	db *gorm.DB
}

func NewPekerjaanAlumniRepository(db *gorm.DB) repository.PekerjaanAlumniRepository {
	return &pekerjaanAlumniRepository{db: db}
}

func (r *pekerjaanAlumniRepository) Create(ctx context.Context, pekerjaan *entity.PekerjaanAlumni) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	query := `INSERT INTO pekerjaan_alumni (mahasiswa_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	result, err := sqlDB.ExecContext(ctx, query,
		pekerjaan.MahasiswaID, pekerjaan.NamaCompany, pekerjaan.Posisi,
		pekerjaan.TanggalMulai, pekerjaan.TanggalSelesai, pekerjaan.Status,
		pekerjaan.Deskripsi, now, now,
	)

	if err != nil {
		return fmt.Errorf("failed to create pekerjaan alumni: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	
	pekerjaan.ID = uint(id)
	pekerjaan.CreatedAt = now
	pekerjaan.UpdatedAt = now
	
	return nil
}

func (r *pekerjaanAlumniRepository) GetByID(ctx context.Context, id uint) (*entity.PekerjaanAlumni, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, mahasiswa_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at 
			  FROM pekerjaan_alumni WHERE id = ? AND deleted_at IS NULL`
	
	pekerjaan := &entity.PekerjaanAlumni{}
	err = sqlDB.QueryRowContext(ctx, query, id).Scan(
		&pekerjaan.ID, &pekerjaan.MahasiswaID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
		&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
		&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get pekerjaan alumni by ID: %w", err)
	}

	return pekerjaan, nil
}

func (r *pekerjaanAlumniRepository) GetByMahasiswaID(ctx context.Context, mahasiswaID uint) ([]*entity.PekerjaanAlumni, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, mahasiswa_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at 
			  FROM pekerjaan_alumni WHERE mahasiswa_id = ? AND deleted_at IS NULL ORDER BY created_at DESC`
	
	rows, err := sqlDB.QueryContext(ctx, query, mahasiswaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pekerjaan alumni by mahasiswa ID: %w", err)
	}
	defer rows.Close()

	var pekerjaanList []*entity.PekerjaanAlumni
	for rows.Next() {
		pekerjaan := &entity.PekerjaanAlumni{}
		err = rows.Scan(
			&pekerjaan.ID, &pekerjaan.MahasiswaID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
			&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
			&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pekerjaan alumni: %w", err)
		}
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}

	return pekerjaanList, nil
}

func (r *pekerjaanAlumniRepository) GetAll(ctx context.Context) ([]*entity.PekerjaanAlumni, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, mahasiswa_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at 
			  FROM pekerjaan_alumni WHERE deleted_at IS NULL ORDER BY created_at DESC`
	
	rows, err := sqlDB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all pekerjaan alumni: %w", err)
	}
	defer rows.Close()

	var pekerjaanList []*entity.PekerjaanAlumni
	for rows.Next() {
		pekerjaan := &entity.PekerjaanAlumni{}
		err = rows.Scan(
			&pekerjaan.ID, &pekerjaan.MahasiswaID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
			&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
			&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pekerjaan alumni: %w", err)
		}
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}

	return pekerjaanList, nil
}

func (r *pekerjaanAlumniRepository) Update(ctx context.Context, pekerjaan *entity.PekerjaanAlumni) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	var setParts []string
	var args []interface{}
	
	if pekerjaan.MahasiswaID != 0 {
		setParts = append(setParts, "mahasiswa_id = ?")
		args = append(args, pekerjaan.MahasiswaID)
	}
	if pekerjaan.NamaCompany != "" {
		setParts = append(setParts, "nama_company = ?")
		args = append(args, pekerjaan.NamaCompany)
	}
	if pekerjaan.Posisi != "" {
		setParts = append(setParts, "posisi = ?")
		args = append(args, pekerjaan.Posisi)
	}
	if !pekerjaan.TanggalMulai.IsZero() {
		setParts = append(setParts, "tanggal_mulai = ?")
		args = append(args, pekerjaan.TanggalMulai)
	}
	if pekerjaan.TanggalSelesai != nil {
		setParts = append(setParts, "tanggal_selesai = ?")
		args = append(args, pekerjaan.TanggalSelesai)
	}
	if pekerjaan.Status != "" {
		setParts = append(setParts, "status = ?")
		args = append(args, pekerjaan.Status)
	}
	if pekerjaan.Deskripsi != "" {
		setParts = append(setParts, "deskripsi = ?")
		args = append(args, pekerjaan.Deskripsi)
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, "updated_at = ?")
	args = append(args, time.Now())
	args = append(args, pekerjaan.ID)

	query := fmt.Sprintf("UPDATE pekerjaan_alumni SET %s WHERE id = ? AND deleted_at IS NULL", strings.Join(setParts, ", "))
	
	result, err := sqlDB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update pekerjaan alumni: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pekerjaan alumni not found or no changes made")
	}

	return nil
}

func (r *pekerjaanAlumniRepository) Delete(ctx context.Context, id uint) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	query := `UPDATE pekerjaan_alumni SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL`
	
	result, err := sqlDB.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete pekerjaan alumni: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pekerjaan alumni not found")
	}

	return nil
}

func (r *pekerjaanAlumniRepository) GetWithPagination(ctx context.Context, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM pekerjaan_alumni WHERE deleted_at IS NULL`
	var total int64
	err = sqlDB.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get paginated results
	query := `SELECT id, mahasiswa_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at 
			  FROM pekerjaan_alumni WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	rows, err := sqlDB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pekerjaan alumni with pagination: %w", err)
	}
	defer rows.Close()

	var pekerjaanList []*entity.PekerjaanAlumni
	for rows.Next() {
		pekerjaan := &entity.PekerjaanAlumni{}
		err = rows.Scan(
			&pekerjaan.ID, &pekerjaan.MahasiswaID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
			&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
			&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pekerjaan alumni: %w", err)
		}
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}

	return pekerjaanList, total, nil
}

func (r *pekerjaanAlumniRepository) GetByIDAndMahasiswaID(ctx context.Context, id, mahasiswaID uint) (*entity.PekerjaanAlumni, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, mahasiswa_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at 
			  FROM pekerjaan_alumni WHERE id = ? AND mahasiswa_id = ? AND deleted_at IS NULL`
	
	pekerjaan := &entity.PekerjaanAlumni{}
	err = sqlDB.QueryRowContext(ctx, query, id, mahasiswaID).Scan(
		&pekerjaan.ID, &pekerjaan.MahasiswaID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
		&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
		&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get pekerjaan alumni by ID and mahasiswa ID: %w", err)
	}

	return pekerjaan, nil
}

func (r *pekerjaanAlumniRepository) GetByMahasiswaIDWithPagination(ctx context.Context, mahasiswaID uint, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM pekerjaan_alumni WHERE mahasiswa_id = ? AND deleted_at IS NULL`
	var total int64
	err = sqlDB.QueryRowContext(ctx, countQuery, mahasiswaID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get paginated results
	query := `SELECT id, mahasiswa_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at 
			  FROM pekerjaan_alumni WHERE mahasiswa_id = ? AND deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	rows, err := sqlDB.QueryContext(ctx, query, mahasiswaID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pekerjaan alumni by mahasiswa ID with pagination: %w", err)
	}
	defer rows.Close()

	var pekerjaanList []*entity.PekerjaanAlumni
	for rows.Next() {
		pekerjaan := &entity.PekerjaanAlumni{}
		err = rows.Scan(
			&pekerjaan.ID, &pekerjaan.MahasiswaID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
			&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
			&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pekerjaan alumni: %w", err)
		}
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}

	return pekerjaanList, total, nil
}

func (r *pekerjaanAlumniRepository) GetWithFilters(ctx context.Context, filters map[string]interface{}) ([]*entity.PekerjaanAlumni, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, mahasiswa_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at 
			  FROM pekerjaan_alumni WHERE deleted_at IS NULL`
	
	var args []interface{}
	var whereClauses []string

	for key, value := range filters {
		switch key {
		case "status":
			whereClauses = append(whereClauses, "status = ?")
			args = append(args, value)
		case "nama_company":
			whereClauses = append(whereClauses, "nama_company LIKE ?")
			args = append(args, fmt.Sprintf("%%%s%%", value))
		case "posisi":
			whereClauses = append(whereClauses, "posisi LIKE ?")
			args = append(args, fmt.Sprintf("%%%s%%", value))
		case "mahasiswa_id":
			whereClauses = append(whereClauses, "mahasiswa_id = ?")
			args = append(args, value)
		}
	}

	if len(whereClauses) > 0 {
		query += " AND " + strings.Join(whereClauses, " AND ")
	}

	query += " ORDER BY created_at DESC"

	rows, err := sqlDB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get pekerjaan alumni with filters: %w", err)
	}
	defer rows.Close()

	var pekerjaanList []*entity.PekerjaanAlumni
	for rows.Next() {
		pekerjaan := &entity.PekerjaanAlumni{}
		err = rows.Scan(
			&pekerjaan.ID, &pekerjaan.MahasiswaID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
			&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
			&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pekerjaan alumni: %w", err)
		}
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}

	return pekerjaanList, nil
}