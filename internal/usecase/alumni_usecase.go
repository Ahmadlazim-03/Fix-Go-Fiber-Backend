package usecase

import (
	"context"
	"errors"
	"strings"

	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/domain/repository"
	"Fix-Go-Fiber-Backend/internal/domain/dto"
)

type AlumniUsecase struct {
	alumniRepo    repository.AlumniRepository
	mahasiswaRepo repository.MahasiswaRepository
}

func NewAlumniUsecase(alumniRepo repository.AlumniRepository, mahasiswaRepo repository.MahasiswaRepository) *AlumniUsecase {
	return &AlumniUsecase{
		alumniRepo:    alumniRepo,
		mahasiswaRepo: mahasiswaRepo,
	}
}

// Implement service.AlumniService interface
func (u *AlumniUsecase) CreateAlumni(ctx context.Context, req *dto.CreateAlumniRequest) (*entity.Alumni, error) {
	// Validate mahasiswa exists
	mahasiswa, err := u.mahasiswaRepo.GetByID(ctx, req.MahasiswaID)
	if err != nil {
		return nil, err
	}
	if mahasiswa == nil {
		return nil, errors.New("mahasiswa tidak ditemukan")
	}

	// Check if alumni already exists for this mahasiswa
	existingAlumni, _ := u.alumniRepo.GetByMahasiswaID(ctx, req.MahasiswaID)
	if existingAlumni != nil {
		return nil, errors.New("mahasiswa sudah terdaftar sebagai alumni")
	}

	alumni := &entity.Alumni{
		MahasiswaID: req.MahasiswaID,
		TahunLulus:  req.TahunLulus,
		NoTelepon:   req.NoTelepon,
		Alamat:      req.Alamat,
	}

	if err := u.alumniRepo.Create(ctx, alumni); err != nil {
		return nil, err
	}

	return alumni, nil
}

func (u *AlumniUsecase) GetAlumniByID(ctx context.Context, id uint) (*entity.Alumni, error) {
	if id == 0 {
		return nil, errors.New("ID tidak valid")
	}

	alumni, err := u.alumniRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if alumni == nil {
		return nil, errors.New("alumni tidak ditemukan")
	}

	return alumni, nil
}

func (u *AlumniUsecase) GetAlumniByMahasiswaID(ctx context.Context, mahasiswaID uint) (*entity.Alumni, error) {
	if mahasiswaID == 0 {
		return nil, errors.New("mahasiswa ID tidak valid")
	}

	return u.alumniRepo.GetByMahasiswaID(ctx, mahasiswaID)
}

func (u *AlumniUsecase) GetAllAlumni(ctx context.Context, search string, limit, offset int) ([]*entity.Alumni, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	if search != "" {
		return u.alumniRepo.Search(ctx, search, limit, offset)
	}

	return u.alumniRepo.GetAll(ctx, limit, offset)
}

func (u *AlumniUsecase) UpdateAlumni(ctx context.Context, id uint, req *dto.UpdateAlumniRequest) (*entity.Alumni, error) {
	if id == 0 {
		return nil, errors.New("ID tidak valid")
	}

	// Check if alumni exists
	existing, err := u.alumniRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("alumni tidak ditemukan")
	}

	// Update only non-zero fields
	if req.TahunLulus > 0 {
		existing.TahunLulus = req.TahunLulus
	}
	if req.NoTelepon != "" {
		existing.NoTelepon = req.NoTelepon
	}
	if req.Alamat != "" {
		existing.Alamat = req.Alamat
	}

	if err := u.alumniRepo.Update(ctx, id, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (u *AlumniUsecase) DeleteAlumni(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("ID tidak valid")
	}

	// Check if alumni exists
	existing, err := u.alumniRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("alumni tidak ditemukan")
	}

	return u.alumniRepo.Delete(ctx, id) // This will soft delete due to gorm.DeletedAt
}

// Legacy methods for backward compatibility
func (u *AlumniUsecase) Create(ctx context.Context, alumni *entity.Alumni) error {
	// Validate required fields
	if err := u.validateAlumni(alumni); err != nil {
		return err
	}

	// Check if mahasiswa exists
	mahasiswa, err := u.mahasiswaRepo.GetByID(ctx, alumni.MahasiswaID)
	if err != nil {
		return err
	}
	if mahasiswa == nil {
		return errors.New("mahasiswa tidak ditemukan")
	}

	// Check if alumni already exists for this mahasiswa
	existingAlumni, _ := u.alumniRepo.GetByMahasiswaID(ctx, alumni.MahasiswaID)
	if existingAlumni != nil {
		return errors.New("mahasiswa sudah terdaftar sebagai alumni")
	}

	return u.alumniRepo.Create(ctx, alumni)
}

func (u *AlumniUsecase) GetByID(ctx context.Context, id uint) (*entity.Alumni, error) {
	if id == 0 {
		return nil, errors.New("ID tidak valid")
	}

	alumni, err := u.alumniRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if alumni == nil {
		return nil, errors.New("alumni tidak ditemukan")
	}

	return alumni, nil
}

func (u *AlumniUsecase) GetWithMahasiswa(ctx context.Context, id uint) (*entity.Alumni, error) {
	if id == 0 {
		return nil, errors.New("ID tidak valid")
	}

	alumni, err := u.alumniRepo.GetWithMahasiswa(ctx, id)
	if err != nil {
		return nil, err
	}
	if alumni == nil {
		return nil, errors.New("alumni tidak ditemukan")
	}

	return alumni, nil
}

func (u *AlumniUsecase) GetAll(ctx context.Context, limit, offset int) ([]*entity.Alumni, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return u.alumniRepo.GetAll(ctx, limit, offset)
}

func (u *AlumniUsecase) Update(ctx context.Context, id uint, alumni *entity.Alumni) error {
	if id == 0 {
		return errors.New("ID tidak valid")
	}

	// Check if alumni exists
	existing, err := u.alumniRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("alumni tidak ditemukan")
	}

	// Validate update data
	if err := u.validateAlumniUpdate(alumni); err != nil {
		return err
	}

	// Check if mahasiswa ID is being changed and if it's valid
	if alumni.MahasiswaID != 0 && alumni.MahasiswaID != existing.MahasiswaID {
		mahasiswa, err := u.mahasiswaRepo.GetByID(ctx, alumni.MahasiswaID)
		if err != nil {
			return err
		}
		if mahasiswa == nil {
			return errors.New("mahasiswa tidak ditemukan")
		}

		// Check if this mahasiswa is already registered as alumni
		existingAlumni, _ := u.alumniRepo.GetByMahasiswaID(ctx, alumni.MahasiswaID)
		if existingAlumni != nil && existingAlumni.ID != id {
			return errors.New("mahasiswa sudah terdaftar sebagai alumni")
		}
	}

	return u.alumniRepo.Update(ctx, id, alumni)
}

func (u *AlumniUsecase) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("ID tidak valid")
	}

	// Check if alumni exists
	existing, err := u.alumniRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("alumni tidak ditemukan")
	}

	return u.alumniRepo.Delete(ctx, id)
}

func (u *AlumniUsecase) Search(ctx context.Context, query string, limit, offset int) ([]*entity.Alumni, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return u.alumniRepo.GetAll(ctx, limit, offset)
	}

	return u.alumniRepo.Search(ctx, query, limit, offset)
}

func (u *AlumniUsecase) validateAlumni(alumni *entity.Alumni) error {
	if alumni.MahasiswaID == 0 {
		return errors.New("mahasiswa ID harus diisi")
	}
	if alumni.TahunLulus <= 0 {
		return errors.New("tahun lulus harus valid")
	}
	return nil
}

func (u *AlumniUsecase) validateAlumniUpdate(alumni *entity.Alumni) error {
	if alumni.TahunLulus < 0 {
		return errors.New("tahun lulus harus valid")
	}
	return nil
}