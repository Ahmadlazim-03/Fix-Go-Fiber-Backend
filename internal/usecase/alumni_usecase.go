package usecase

import (
	"context"
	"errors"
	"strings"

	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/domain/repository"
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