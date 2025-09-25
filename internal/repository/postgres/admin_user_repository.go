package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	query := `INSERT INTO admin_users (username, email, password, role, is_active, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`
	
	now := time.Now()
	err = sqlDB.QueryRowContext(ctx, query,
		admin.Username, admin.Email, admin.Password, admin.Role,
		admin.IsActive, now, now,
	).Scan(&admin.ID)

	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}
	
	admin.CreatedAt = now
	admin.UpdatedAt = now
	return nil
}

func (r *adminUserRepository) GetByID(ctx context.Context, id uint) (*entity.AdminUser, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, username, email, password, role, is_active, created_at, updated_at 
			  FROM admin_users WHERE id = ? AND deleted_at IS NULL`
	
	var admin entity.AdminUser
	err = sqlDB.QueryRowContext(ctx, query, id).Scan(
		&admin.ID, &admin.Username, &admin.Email, &admin.Password,
		&admin.Role, &admin.IsActive, &admin.CreatedAt, &admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get admin user by ID: %w", err)
	}

	return &admin, nil
}

func (r *adminUserRepository) GetByUsername(ctx context.Context, username string) (*entity.AdminUser, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, username, email, password, role, is_active, created_at, updated_at 
			  FROM admin_users WHERE username = ? AND deleted_at IS NULL`
	
	var admin entity.AdminUser
	err = sqlDB.QueryRowContext(ctx, query, username).Scan(
		&admin.ID, &admin.Username, &admin.Email, &admin.Password,
		&admin.Role, &admin.IsActive, &admin.CreatedAt, &admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get admin user by username: %w", err)
	}

	return &admin, nil
}

func (r *adminUserRepository) GetByEmail(ctx context.Context, email string) (*entity.AdminUser, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, username, email, password, role, is_active, created_at, updated_at 
			  FROM admin_users WHERE email = ? AND deleted_at IS NULL`
	
	var admin entity.AdminUser
	err = sqlDB.QueryRowContext(ctx, query, email).Scan(
		&admin.ID, &admin.Username, &admin.Email, &admin.Password,
		&admin.Role, &admin.IsActive, &admin.CreatedAt, &admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get admin user by email: %w", err)
	}

	return &admin, nil
}

func (r *adminUserRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.AdminUser, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM admin_users WHERE deleted_at IS NULL`
	var total int64
	err = sqlDB.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count admin users: %w", err)
	}

	// Get data with pagination
	query := `SELECT id, username, email, password, role, is_active, created_at, updated_at 
			  FROM admin_users WHERE deleted_at IS NULL 
			  ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	rows, err := sqlDB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get admin users list: %w", err)
	}
	defer rows.Close()

	var admins []*entity.AdminUser
	for rows.Next() {
		var admin entity.AdminUser
		err = rows.Scan(
			&admin.ID, &admin.Username, &admin.Email, &admin.Password,
			&admin.Role, &admin.IsActive, &admin.CreatedAt, &admin.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan admin user: %w", err)
		}
		admins = append(admins, &admin)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating admin users rows: %w", err)
	}

	return admins, total, nil
}

func (r *adminUserRepository) Update(ctx context.Context, id uint, admin *entity.AdminUser) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	query := `UPDATE admin_users SET username = ?, email = ?, role = ?, is_active = ?, updated_at = ? 
			  WHERE id = ? AND deleted_at IS NULL`
	
	result, err := sqlDB.ExecContext(ctx, query, 
		admin.Username, admin.Email, admin.Role, admin.IsActive, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update admin user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("admin user not found or already deleted")
	}

	return nil
}

func (r *adminUserRepository) Delete(ctx context.Context, id uint) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	query := `UPDATE admin_users SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL`
	
	result, err := sqlDB.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete admin user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("admin user not found or already deleted")
	}

	return nil
}

func (r *adminUserRepository) GetActiveAdmins(ctx context.Context, limit, offset int) ([]*entity.AdminUser, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	// Count total active admins
	countQuery := `SELECT COUNT(*) FROM admin_users WHERE is_active = true AND deleted_at IS NULL`
	var total int64
	err = sqlDB.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count active admin users: %w", err)
	}

	// Get active admins with pagination
	query := `SELECT id, username, email, password, role, is_active, created_at, updated_at 
			  FROM admin_users WHERE is_active = true AND deleted_at IS NULL 
			  ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	rows, err := sqlDB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get active admin users: %w", err)
	}
	defer rows.Close()

	var admins []*entity.AdminUser
	for rows.Next() {
		var admin entity.AdminUser
		err = rows.Scan(
			&admin.ID, &admin.Username, &admin.Email, &admin.Password,
			&admin.Role, &admin.IsActive, &admin.CreatedAt, &admin.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan active admin user: %w", err)
		}
		admins = append(admins, &admin)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating active admin users: %w", err)
	}

	return admins, total, nil
}
