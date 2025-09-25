package service

import (
	"context"
	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/domain/dto"
)

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Username string `json:"username,omitempty"` // for admin
}

// AuthService interface untuk authentication domain services
type AuthService interface {
	// Login methods
	LoginMahasiswa(req *dto.MahasiswaLoginRequest) (*dto.LoginResponse, error)
	LoginAlumni(req *dto.AlumniLoginRequest) (*dto.LoginResponse, error)
	LoginAdmin(req *dto.AdminLoginRequest) (*dto.LoginResponse, error)
	
	// Register methods
	RegisterMahasiswa(req *dto.RegisterMahasiswaRequest) (*dto.RegisterResponse, error)
	RegisterAlumni(req *dto.RegisterAlumniRequest) (*dto.RegisterResponse, error)
	
	// Token validation
	ValidateToken(token string) (*JWTClaims, error)
	
	// Legacy methods for backward compatibility
	ValidateCredentials(ctx context.Context, email, password string) (*entity.Mahasiswa, error)
	ValidateAdminCredentials(ctx context.Context, username, password string) (*entity.AdminUser, error)
	GenerateToken(ctx context.Context, user interface{}) (string, error)
}

// EmailService interface untuk email domain services
type EmailService interface {
	SendWelcomeEmail(ctx context.Context, email, name string) error
	SendPasswordResetEmail(ctx context.Context, email, resetToken string) error
	SendAlumniRegistrationNotification(ctx context.Context, alumni *entity.Alumni) error
}

// NotificationService interface untuk notification domain services
type NotificationService interface {
	SendNotification(ctx context.Context, userID uint, message string) error
	SendBulkNotification(ctx context.Context, userIDs []uint, message string) error
}