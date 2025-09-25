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
	alumniRepo    repository.AlumniRepository
	adminRepo     repository.AdminUserRepository
	jwtUtil       *jwt.JWTUtil
	bcryptUtil    *bcrypt.BcryptUtil
}

func NewAuthService(
	mahasiswaRepo repository.MahasiswaRepository,
	alumniRepo repository.AlumniRepository,
	adminRepo repository.AdminUserRepository,
	jwtUtil *jwt.JWTUtil,
	bcryptUtil *bcrypt.BcryptUtil,
) service.AuthService {
	return &authService{
		mahasiswaRepo: mahasiswaRepo,
		alumniRepo:    alumniRepo,
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
	
	// Find alumni by email through mahasiswa relation
	mahasiswa, err := s.mahasiswaRepo.GetByEmail(ctx, req.Email)
	if err != nil || mahasiswa == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if !s.bcryptUtil.CheckPasswordHash(req.Password, mahasiswa.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Find alumni record by mahasiswa_id
	alumni, err := s.alumniRepo.GetByMahasiswaID(ctx, mahasiswa.ID)
	if err != nil {
		return nil, errors.New("alumni record not found")
	}

	// Generate JWT token
	claims := &service.JWTClaims{
		UserID: alumni.ID,
		Email:  mahasiswa.Email,
		Role:   "alumni",
	}

	token, expiresAt, err := s.jwtUtil.GenerateToken(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token:     token,
		User:      alumni.ToResponse(),
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

// RegisterAlumni creates a new alumni account
func (s *authService) RegisterAlumni(req *dto.RegisterAlumniRequest) (*dto.RegisterResponse, error) {
	ctx := context.Background()
	
	// First check if a mahasiswa with this NIM exists
	mahasiswa, err := s.mahasiswaRepo.GetByNIM(ctx, req.NIM)
	if err != nil || mahasiswa == nil {
		return nil, errors.New("mahasiswa with this NIM not found. Please register as mahasiswa first")
	}
	
	// Check if this mahasiswa is already registered as alumni
	existingAlumni, _ := s.alumniRepo.GetByMahasiswaID(ctx, mahasiswa.ID)
	if existingAlumni != nil {
		return nil, errors.New("alumni with this NIM already exists")
	}
	
	// Verify if the provided details match the mahasiswa record
	if mahasiswa.Nama != req.Nama || mahasiswa.Email != req.Email || mahasiswa.Jurusan != req.Jurusan {
		return nil, errors.New("provided details do not match mahasiswa record")
	}
	
	// Create alumni entity
	alumni := &entity.Alumni{
		MahasiswaID: mahasiswa.ID,
		TahunLulus:  req.TahunLulus,
	}
	
	// Save to database
	if err := s.alumniRepo.Create(ctx, alumni); err != nil {
		return nil, fmt.Errorf("failed to create alumni: %w", err)
	}
	
	return &dto.RegisterResponse{
		ID:      int64(alumni.ID),
		Message: "Alumni registered successfully",
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