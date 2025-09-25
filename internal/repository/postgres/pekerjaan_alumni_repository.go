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

type pekerjaanAlumniRepository struct {
	db *gorm.DB
}

func NewPekerjaanAlumniRepository(db *gorm.DB) repository.PekerjaanAlumniRepository {
	return &pekerjaanAlumniRepository{
		db: db,
	}
}

func (r *pekerjaanAlumniRepository) Create(ctx context.Context, pekerjaan *entity.PekerjaanAlumni) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	query := `INSERT INTO pekerjaan_alumni (alumni_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	
	now := time.Now()
	err = sqlDB.QueryRowContext(ctx, query,
		pekerjaan.AlumniID, pekerjaan.NamaCompany, pekerjaan.Posisi,
		pekerjaan.TanggalMulai, pekerjaan.TanggalSelesai, pekerjaan.Status,
		pekerjaan.Deskripsi, now, now,
	).Scan(&pekerjaan.ID)

	if err != nil {
		return fmt.Errorf("failed to create pekerjaan alumni: %w", err)
	}
	
	pekerjaan.CreatedAt = now
	pekerjaan.UpdatedAt = now
	return nil
}

func (r *pekerjaanAlumniRepository) GetByID(ctx context.Context, id uint) (*entity.PekerjaanAlumni, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, alumni_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at 
			  FROM pekerjaan_alumni WHERE id = $1 AND deleted_at IS NULL`
	
	var pekerjaan entity.PekerjaanAlumni
	err = sqlDB.QueryRowContext(ctx, query, id).Scan(
		&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
		&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
		&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get pekerjaan alumni by ID: %w", err)
	}

	return &pekerjaan, nil
}

func (r *pekerjaanAlumniRepository) GetByAlumniID(ctx context.Context, alumniID uint) ([]*entity.PekerjaanAlumni, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, alumni_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at 
			  FROM pekerjaan_alumni WHERE alumni_id = $1 AND deleted_at IS NULL
			  ORDER BY created_at DESC`
	
	rows, err := sqlDB.QueryContext(ctx, query, alumniID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pekerjaan by alumni ID: %w", err)
	}
	defer rows.Close()

	var pekerjaans []*entity.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan entity.PekerjaanAlumni
		err = rows.Scan(
			&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
			&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
			&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pekerjaan alumni: %w", err)
		}
		pekerjaans = append(pekerjaans, &pekerjaan)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating pekerjaan rows: %w", err)
	}

	return pekerjaans, nil
}

func (r *pekerjaanAlumniRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM pekerjaan_alumni WHERE deleted_at IS NULL`
	var total int64
	err = sqlDB.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pekerjaan alumni: %w", err)
	}

	// Get data with pagination
	query := `SELECT id, alumni_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at 
			  FROM pekerjaan_alumni WHERE deleted_at IS NULL 
			  ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	rows, err := sqlDB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pekerjaan alumni list: %w", err)
	}
	defer rows.Close()

	var pekerjaans []*entity.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan entity.PekerjaanAlumni
		err = rows.Scan(
			&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
			&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
			&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pekerjaan alumni: %w", err)
		}
		pekerjaans = append(pekerjaans, &pekerjaan)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating pekerjaan rows: %w", err)
	}

	return pekerjaans, total, nil
}

func (r *pekerjaanAlumniRepository) Update(ctx context.Context, id uint, pekerjaan *entity.PekerjaanAlumni) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if pekerjaan.AlumniID != 0 {
		setParts = append(setParts, fmt.Sprintf("alumni_id = $%d", argIndex))
		args = append(args, pekerjaan.AlumniID)
		argIndex++
	}
	if pekerjaan.NamaCompany != "" {
		setParts = append(setParts, fmt.Sprintf("nama_company = $%d", argIndex))
		args = append(args, pekerjaan.NamaCompany)
		argIndex++
	}
	if pekerjaan.Posisi != "" {
		setParts = append(setParts, fmt.Sprintf("posisi = $%d", argIndex))
		args = append(args, pekerjaan.Posisi)
		argIndex++
	}
	if !pekerjaan.TanggalMulai.IsZero() {
		setParts = append(setParts, fmt.Sprintf("tanggal_mulai = $%d", argIndex))
		args = append(args, pekerjaan.TanggalMulai)
		argIndex++
	}
	if pekerjaan.TanggalSelesai != nil {
		setParts = append(setParts, fmt.Sprintf("tanggal_selesai = $%d", argIndex))
		args = append(args, pekerjaan.TanggalSelesai)
		argIndex++
	}
	if pekerjaan.Status != "" {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, pekerjaan.Status)
		argIndex++
	}
	if pekerjaan.Deskripsi != "" {
		setParts = append(setParts, fmt.Sprintf("deskripsi = $%d", argIndex))
		args = append(args, pekerjaan.Deskripsi)
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

	query := fmt.Sprintf("UPDATE pekerjaan_alumni SET %s WHERE id = $%d AND deleted_at IS NULL", 
		strings.Join(setParts, ", "), argIndex)

	result, err := sqlDB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update pekerjaan alumni: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("pekerjaan alumni not found or already deleted")
	}

	return nil
}

func (r *pekerjaanAlumniRepository) Delete(ctx context.Context, id uint) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	query := `UPDATE pekerjaan_alumni SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	
	result, err := sqlDB.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete pekerjaan alumni: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("pekerjaan alumni not found or already deleted")
	}

	return nil
}

func (r *pekerjaanAlumniRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	searchQuery := "%" + strings.ToLower(query) + "%"

	// Count total
	countSQL := `SELECT COUNT(*) FROM pekerjaan_alumni 
				 WHERE deleted_at IS NULL AND (
					 LOWER(nama_company) LIKE $1 OR LOWER(posisi) LIKE $1 OR 
					 LOWER(status) LIKE $1 OR LOWER(deskripsi) LIKE $1
				 )`
	var total int64
	err = sqlDB.QueryRowContext(ctx, countSQL, searchQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Get data with pagination
	dataSQL := `SELECT id, alumni_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at 
				FROM pekerjaan_alumni 
				WHERE deleted_at IS NULL AND (
					LOWER(nama_company) LIKE $1 OR LOWER(posisi) LIKE $1 OR 
					LOWER(status) LIKE $1 OR LOWER(deskripsi) LIKE $1
				)
				ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := sqlDB.QueryContext(ctx, dataSQL, searchQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search pekerjaan alumni: %w", err)
	}
	defer rows.Close()

	var pekerjaans []*entity.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan entity.PekerjaanAlumni
		err = rows.Scan(
			&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
			&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
			&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pekerjaan alumni: %w", err)
		}
		pekerjaans = append(pekerjaans, &pekerjaan)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating search results: %w", err)
	}

	return pekerjaans, total, nil
}

func (r *pekerjaanAlumniRepository) GetByIDAndAlumniID(ctx context.Context, id, alumniID uint) (*entity.PekerjaanAlumni, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, alumni_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at 
			  FROM pekerjaan_alumni WHERE id = $1 AND alumni_id = $2 AND deleted_at IS NULL`
	
	var pekerjaan entity.PekerjaanAlumni
	err = sqlDB.QueryRowContext(ctx, query, id, alumniID).Scan(
		&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
		&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
		&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get pekerjaan alumni by ID and alumni ID: %w", err)
	}

	return &pekerjaan, nil
}

func (r *pekerjaanAlumniRepository) GetByAlumniIDWithPagination(ctx context.Context, alumniID uint, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM pekerjaan_alumni WHERE alumni_id = $1 AND deleted_at IS NULL`
	var total int64
	err = sqlDB.QueryRowContext(ctx, countQuery, alumniID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pekerjaan alumni: %w", err)
	}

	// Get data with pagination
	query := `SELECT id, alumni_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at 
			  FROM pekerjaan_alumni WHERE alumni_id = $1 AND deleted_at IS NULL 
			  ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	
	rows, err := sqlDB.QueryContext(ctx, query, alumniID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pekerjaan alumni list: %w", err)
	}
	defer rows.Close()

	var pekerjaans []*entity.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan entity.PekerjaanAlumni
		err = rows.Scan(
			&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
			&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
			&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pekerjaan alumni: %w", err)
		}
		pekerjaans = append(pekerjaans, &pekerjaan)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating pekerjaan rows: %w", err)
	}

	return pekerjaans, total, nil
}

func (r *pekerjaanAlumniRepository) GetWithAlumni(ctx context.Context, id uint) (*entity.PekerjaanAlumni, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT p.id, p.alumni_id, p.nama_company, p.posisi, p.tanggal_mulai, p.tanggal_selesai, 
					 p.status, p.deskripsi, p.created_at, p.updated_at,
					 a.id, a.mahasiswa_id, a.tahun_lulus, a.no_telepon, a.alamat, a.created_at, a.updated_at
			  FROM pekerjaan_alumni p
			  JOIN alumni a ON p.alumni_id = a.id
			  WHERE p.id = $1 AND p.deleted_at IS NULL AND a.deleted_at IS NULL`
	
	var pekerjaan entity.PekerjaanAlumni
	var alumni entity.Alumni
	
	err = sqlDB.QueryRowContext(ctx, query, id).Scan(
		&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
		&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
		&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		&alumni.ID, &alumni.MahasiswaID, &alumni.TahunLulus, &alumni.NoTelepon,
		&alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get pekerjaan alumni with alumni: %w", err)
	}

	pekerjaan.Alumni = alumni
	return &pekerjaan, nil
}

func (r *pekerjaanAlumniRepository) GetActiveJobs(ctx context.Context, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	// Count total active jobs
	countQuery := `SELECT COUNT(*) FROM pekerjaan_alumni WHERE status = 'aktif' AND deleted_at IS NULL`
	var total int64
	err = sqlDB.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count active jobs: %w", err)
	}

	// Get active jobs with pagination
	query := `SELECT id, alumni_id, nama_company, posisi, tanggal_mulai, tanggal_selesai, status, deskripsi, created_at, updated_at 
			  FROM pekerjaan_alumni WHERE status = 'aktif' AND deleted_at IS NULL 
			  ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	rows, err := sqlDB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get active jobs: %w", err)
	}
	defer rows.Close()

	var pekerjaans []*entity.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan entity.PekerjaanAlumni
		err = rows.Scan(
			&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaCompany, &pekerjaan.Posisi,
			&pekerjaan.TanggalMulai, &pekerjaan.TanggalSelesai, &pekerjaan.Status,
			&pekerjaan.Deskripsi, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan active job: %w", err)
		}
		pekerjaans = append(pekerjaans, &pekerjaan)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating active jobs: %w", err)
	}

	return pekerjaans, total, nil
}