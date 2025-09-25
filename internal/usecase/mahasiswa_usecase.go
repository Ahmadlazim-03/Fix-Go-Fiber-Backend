package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/domain/repository"
	"Fix-Go-Fiber-Backend/pkg/bcrypt"
)

type MahasiswaUsecase struct {
	mahasiswaRepo repository.MahasiswaRepository
	bcryptHelper  bcrypt.BcryptHelper
}

func NewMahasiswaUsecase(mahasiswaRepo repository.MahasiswaRepository, bcryptHelper bcrypt.BcryptHelper) *MahasiswaUsecase {
	return &MahasiswaUsecase{
		mahasiswaRepo: mahasiswaRepo,
		bcryptHelper:  bcryptHelper,
	}
}

func (u *MahasiswaUsecase) Create(ctx context.Context, mahasiswa *entity.Mahasiswa) error {
	// Validate required fields
	if err := u.validateMahasiswa(mahasiswa); err != nil {
		return err
	}

	// Check if NIM already exists
	existingByNIM, _ := u.mahasiswaRepo.GetByNIM(ctx, mahasiswa.NIM)
	if existingByNIM != nil {
		return errors.New("NIM sudah terdaftar")
	}

	// Check if email already exists
	existingByEmail, _ := u.mahasiswaRepo.GetByEmail(ctx, mahasiswa.Email)
	if existingByEmail != nil {
		return errors.New("email sudah terdaftar")
	}

	// Hash password
	hashedPassword, err := u.bcryptHelper.HashPassword(mahasiswa.Password)
	if err != nil {
		return fmt.Errorf("gagal hash password: %w", err)
	}
	mahasiswa.Password = hashedPassword

	return u.mahasiswaRepo.Create(ctx, mahasiswa)
}

func (u *MahasiswaUsecase) GetByID(ctx context.Context, id uint) (*entity.Mahasiswa, error) {
	if id == 0 {
		return nil, errors.New("ID tidak valid")
	}

	mahasiswa, err := u.mahasiswaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if mahasiswa == nil {
		return nil, errors.New("mahasiswa tidak ditemukan")
	}

	return mahasiswa, nil
}

func (u *MahasiswaUsecase) GetAll(ctx context.Context, limit, offset int) ([]*entity.Mahasiswa, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return u.mahasiswaRepo.GetAll(ctx, limit, offset)
}

func (u *MahasiswaUsecase) Update(ctx context.Context, id uint, mahasiswa *entity.Mahasiswa) error {
	if id == 0 {
		return errors.New("ID tidak valid")
	}

	// Check if mahasiswa exists
	existing, err := u.mahasiswaRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("mahasiswa tidak ditemukan")
	}

	// Validate update data
	if err := u.validateMahasiswaUpdate(mahasiswa); err != nil {
		return err
	}

	// Check NIM uniqueness if changed
	if mahasiswa.NIM != "" && mahasiswa.NIM != existing.NIM {
		existingByNIM, _ := u.mahasiswaRepo.GetByNIM(ctx, mahasiswa.NIM)
		if existingByNIM != nil && existingByNIM.ID != id {
			return errors.New("NIM sudah terdaftar")
		}
	}

	// Check email uniqueness if changed
	if mahasiswa.Email != "" && mahasiswa.Email != existing.Email {
		existingByEmail, _ := u.mahasiswaRepo.GetByEmail(ctx, mahasiswa.Email)
		if existingByEmail != nil && existingByEmail.ID != id {
			return errors.New("email sudah terdaftar")
		}
	}

	// Hash password if provided
	if mahasiswa.Password != "" {
		hashedPassword, err := u.bcryptHelper.HashPassword(mahasiswa.Password)
		if err != nil {
			return fmt.Errorf("gagal hash password: %w", err)
		}
		mahasiswa.Password = hashedPassword
	}

	return u.mahasiswaRepo.Update(ctx, id, mahasiswa)
}

func (u *MahasiswaUsecase) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("ID tidak valid")
	}

	// Check if mahasiswa exists
	existing, err := u.mahasiswaRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("mahasiswa tidak ditemukan")
	}

	return u.mahasiswaRepo.Delete(ctx, id)
}

func (u *MahasiswaUsecase) Search(ctx context.Context, query string, limit, offset int) ([]*entity.Mahasiswa, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return u.mahasiswaRepo.GetAll(ctx, limit, offset)
	}

	return u.mahasiswaRepo.Search(ctx, query, limit, offset)
}

func (u *MahasiswaUsecase) validateMahasiswa(mahasiswa *entity.Mahasiswa) error {
	if mahasiswa.NIM == "" {
		return errors.New("NIM harus diisi")
	}
	if mahasiswa.Nama == "" {
		return errors.New("nama harus diisi")
	}
	if mahasiswa.Jurusan == "" {
		return errors.New("jurusan harus diisi")
	}
	if mahasiswa.Angkatan <= 0 {
		return errors.New("angkatan harus valid")
	}
	if mahasiswa.Email == "" {
		return errors.New("email harus diisi")
	}
	if mahasiswa.Password == "" {
		return errors.New("password harus diisi")
	}
	return nil
}

func (u *MahasiswaUsecase) validateMahasiswaUpdate(mahasiswa *entity.Mahasiswa) error {
	if mahasiswa.NIM != "" && len(strings.TrimSpace(mahasiswa.NIM)) == 0 {
		return errors.New("NIM tidak boleh kosong")
	}
	if mahasiswa.Nama != "" && len(strings.TrimSpace(mahasiswa.Nama)) == 0 {
		return errors.New("nama tidak boleh kosong")
	}
	if mahasiswa.Jurusan != "" && len(strings.TrimSpace(mahasiswa.Jurusan)) == 0 {
		return errors.New("jurusan tidak boleh kosong")
	}
	if mahasiswa.Email != "" && len(strings.TrimSpace(mahasiswa.Email)) == 0 {
		return errors.New("email tidak boleh kosong")
	}
	return nil
}