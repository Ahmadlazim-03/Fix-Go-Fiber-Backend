package usecase

import (
	"context"
	"errors"
	"fmt"

	"Fix-Go-Fiber-Backend/internal/domain/dto"
	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/domain/repository"
	"Fix-Go-Fiber-Backend/internal/domain/service"
	"Fix-Go-Fiber-Backend/pkg/bcrypt"
	"Fix-Go-Fiber-Backend/pkg/jwt"
)

type authService struct {
	mahasiswaRepo repository.MahasiswaRepository
	adminRepo     repository.AdminUserRepository
	jwtUtil       *jwt.JWTUtil
	bcryptUtil    *bcrypt.BcryptUtil
}

func NewAuthService(
	mahasiswaRepo repository.MahasiswaRepository,
	adminRepo repository.AdminUserRepository,
	jwtUtil *jwt.JWTUtil,
	bcryptUtil *bcrypt.BcryptUtil,
) service.AuthService {
	return &authService{
		mahasiswaRepo: mahasiswaRepo,
		adminRepo:     adminRepo,
		jwtUtil:       jwtUtil,
		bcryptUtil:    bcryptUtil,
	}
}

func (s *authService) LoginMahasiswa(req *dto.MahasiswaLoginRequest) (*dto.LoginResponse, error) {
	ctx := context.Background()
	
	// Find mahasiswa by email
	mahasiswa, err := s.mahasiswaRepo.GetByEmail(ctx, req.Email)
	if err != nil || mahasiswa == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if !s.bcryptUtil.CheckPasswordHash(req.Password, mahasiswa.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	claims := &service.JWTClaims{
		UserID: mahasiswa.ID,
		Email:  mahasiswa.Email,
		Role:   "mahasiswa",
	}

	token, expiresAt, err := s.jwtUtil.GenerateToken(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token:     token,
		User:      mahasiswa.ToResponse(),
		Role:      "mahasiswa",
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

func (s *authService) LoginAlumni(req *dto.AlumniLoginRequest) (*dto.LoginResponse, error) {
	ctx := context.Background()
	
	// Find mahasiswa by email
	mahasiswa, err := s.mahasiswaRepo.GetByEmail(ctx, req.Email)
	if err != nil || mahasiswa == nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if mahasiswa has graduated (is alumni)
	if !mahasiswa.IsAlumni() {
		return nil, errors.New("account is not an alumni account")
	}

	// Verify password
	if !s.bcryptUtil.CheckPasswordHash(req.Password, mahasiswa.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	claims := &service.JWTClaims{
		UserID: mahasiswa.ID,
		Email:  mahasiswa.Email,
		Role:   "alumni",
	}

	token, expiresAt, err := s.jwtUtil.GenerateToken(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token:     token,
		User:      mahasiswa.ToResponse(),
		Role:      "alumni",
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

// RegisterMahasiswa creates a new mahasiswa account
func (s *authService) RegisterMahasiswa(req *dto.RegisterMahasiswaRequest) (*dto.RegisterResponse, error) {
	ctx := context.Background()
	
	// Check if email already exists
	existingMahasiswa, _ := s.mahasiswaRepo.GetByEmail(ctx, req.Email)
	if existingMahasiswa != nil {
		return nil, errors.New("email already registered")
	}
	
	// Check if NIM already exists
	existingByNIM, _ := s.mahasiswaRepo.GetByNIM(ctx, req.NIM)
	if existingByNIM != nil {
		return nil, errors.New("NIM already registered")
	}
	
	// Hash password
	hashedPassword, err := s.bcryptUtil.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Create mahasiswa entity
	mahasiswa := &entity.Mahasiswa{
		NIM:      req.NIM,
		Nama:     req.Nama,
		Email:    req.Email,
		Password: hashedPassword,
		Jurusan:  req.Jurusan,
		Angkatan: req.Angkatan,
	}
	
	// Save to database
	if err := s.mahasiswaRepo.Create(ctx, mahasiswa); err != nil {
		return nil, fmt.Errorf("failed to create mahasiswa: %w", err)
	}
	
	return &dto.RegisterResponse{
		ID:      int64(mahasiswa.ID),
		Message: "Mahasiswa registered successfully",
	}, nil
}

// GraduateMahasiswa marks a mahasiswa as graduated (alumni)
func (s *authService) GraduateMahasiswa(req *dto.GraduateMahasiswaRequest) (*dto.RegisterResponse, error) {
	ctx := context.Background()
	
	// Find mahasiswa by ID
	mahasiswa, err := s.mahasiswaRepo.GetByID(ctx, req.MahasiswaID)
	if err != nil || mahasiswa == nil {
		return nil, errors.New("mahasiswa not found")
	}
	
	// Check if already graduated
	if mahasiswa.IsAlumni() {
		return nil, errors.New("mahasiswa is already graduated")
	}
	
	// Graduate the mahasiswa
	mahasiswa.Graduate(req.TahunLulus, req.NoTelepon, req.AlamatAlumni)
	
	// Update in database
	if err := s.mahasiswaRepo.Update(ctx, req.MahasiswaID, mahasiswa); err != nil {
		return nil, fmt.Errorf("failed to graduate mahasiswa: %w", err)
	}
	
	return &dto.RegisterResponse{
		ID:      int64(mahasiswa.ID),
		Message: "Mahasiswa graduated successfully",
	}, nil
}

func (s *authService) LoginAdmin(req *dto.AdminLoginRequest) (*dto.LoginResponse, error) {
	ctx := context.Background()
	
	// Find admin by username
	admin, err := s.adminRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if admin is active
	if !admin.IsActive {
		return nil, errors.New("admin account is inactive")
	}

	// Verify password
	if !s.bcryptUtil.CheckPasswordHash(req.Password, admin.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	claims := &service.JWTClaims{
		UserID:   admin.ID,
		Email:    admin.Email,
		Role:     "admin",
		Username: admin.Username,
	}

	token, expiresAt, err := s.jwtUtil.GenerateToken(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token:     token,
		User:      admin.ToResponse(),
		Role:      "admin",
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

func (s *authService) ValidateToken(token string) (*service.JWTClaims, error) {
	return s.jwtUtil.ValidateToken(token)
}

// Legacy methods for backward compatibility
func (s *authService) ValidateCredentials(ctx context.Context, email, password string) (*entity.Mahasiswa, error) {
	mahasiswa, err := s.mahasiswaRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !s.bcryptUtil.CheckPasswordHash(password, mahasiswa.Password) {
		return nil, errors.New("invalid password")
	}

	return mahasiswa, nil
}

func (s *authService) ValidateAdminCredentials(ctx context.Context, username, password string) (*entity.AdminUser, error) {
	admin, err := s.adminRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if !admin.IsActive {
		return nil, errors.New("admin account is inactive")
	}

	if !s.bcryptUtil.CheckPasswordHash(password, admin.Password) {
		return nil, errors.New("invalid password")
	}

	return admin, nil
}

func (s *authService) GenerateToken(ctx context.Context, user interface{}) (string, error) {
	var claims *service.JWTClaims

	switch u := user.(type) {
	case *entity.Mahasiswa:
		claims = &service.JWTClaims{
			UserID: u.ID,
			Email:  u.Email,
			Role:   "mahasiswa",
		}
	case *entity.AdminUser:
		claims = &service.JWTClaims{
			UserID:   u.ID,
			Email:    u.Email,
			Role:     "admin",
			Username: u.Username,
		}
	default:
		return "", errors.New("unsupported user type")
	}

	token, _, err := s.jwtUtil.GenerateToken(claims)
	return token, err
}