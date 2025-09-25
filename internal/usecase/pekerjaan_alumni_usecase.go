package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/domain/repository"
)

type PekerjaanAlumniUsecase struct {
	pekerjaanRepo repository.PekerjaanAlumniRepository
	alumniRepo    repository.AlumniRepository
}

func NewPekerjaanAlumniUsecase(pekerjaanRepo repository.PekerjaanAlumniRepository, alumniRepo repository.AlumniRepository) *PekerjaanAlumniUsecase {
	return &PekerjaanAlumniUsecase{
		pekerjaanRepo: pekerjaanRepo,
		alumniRepo:    alumniRepo,
	}
}

func (u *PekerjaanAlumniUsecase) Create(ctx context.Context, pekerjaan *entity.PekerjaanAlumni) error {
	// Validate required fields
	if err := u.validatePekerjaan(pekerjaan); err != nil {
		return err
	}

	// Check if alumni exists
	alumni, err := u.alumniRepo.GetByID(ctx, pekerjaan.AlumniID)
	if err != nil {
		return err
	}
	if alumni == nil {
		return errors.New("alumni tidak ditemukan")
	}

	// Set default status if empty
	if pekerjaan.Status == "" {
		pekerjaan.Status = entity.StatusAktif
	}

	return u.pekerjaanRepo.Create(ctx, pekerjaan)
}

func (u *PekerjaanAlumniUsecase) GetByID(ctx context.Context, id uint) (*entity.PekerjaanAlumni, error) {
	if id == 0 {
		return nil, errors.New("ID tidak valid")
	}

	pekerjaan, err := u.pekerjaanRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if pekerjaan == nil {
		return nil, errors.New("pekerjaan tidak ditemukan")
	}

	return pekerjaan, nil
}

func (u *PekerjaanAlumniUsecase) GetWithAlumni(ctx context.Context, id uint) (*entity.PekerjaanAlumni, error) {
	if id == 0 {
		return nil, errors.New("ID tidak valid")
	}

	pekerjaan, err := u.pekerjaanRepo.GetWithAlumni(ctx, id)
	if err != nil {
		return nil, err
	}
	if pekerjaan == nil {
		return nil, errors.New("pekerjaan tidak ditemukan")
	}

	return pekerjaan, nil
}

func (u *PekerjaanAlumniUsecase) GetByAlumniID(ctx context.Context, alumniID uint) ([]*entity.PekerjaanAlumni, error) {
	if alumniID == 0 {
		return nil, errors.New("alumni ID tidak valid")
	}

	// Check if alumni exists
	alumni, err := u.alumniRepo.GetByID(ctx, alumniID)
	if err != nil {
		return nil, err
	}
	if alumni == nil {
		return nil, errors.New("alumni tidak ditemukan")
	}

	return u.pekerjaanRepo.GetByAlumniID(ctx, alumniID)
}

func (u *PekerjaanAlumniUsecase) GetAll(ctx context.Context, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return u.pekerjaanRepo.GetAll(ctx, limit, offset)
}

func (u *PekerjaanAlumniUsecase) GetActiveJobs(ctx context.Context, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return u.pekerjaanRepo.GetActiveJobs(ctx, limit, offset)
}

func (u *PekerjaanAlumniUsecase) Update(ctx context.Context, id uint, pekerjaan *entity.PekerjaanAlumni) error {
	if id == 0 {
		return errors.New("ID tidak valid")
	}

	// Check if pekerjaan exists
	existing, err := u.pekerjaanRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("pekerjaan tidak ditemukan")
	}

	// Validate update data
	if err := u.validatePekerjaanUpdate(pekerjaan); err != nil {
		return err
	}

	// Check if alumni ID is being changed and if it's valid
	if pekerjaan.AlumniID != 0 && pekerjaan.AlumniID != existing.AlumniID {
		alumni, err := u.alumniRepo.GetByID(ctx, pekerjaan.AlumniID)
		if err != nil {
			return err
		}
		if alumni == nil {
			return errors.New("alumni tidak ditemukan")
		}
	}

	return u.pekerjaanRepo.Update(ctx, id, pekerjaan)
}

func (u *PekerjaanAlumniUsecase) CompleteJob(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("ID tidak valid")
	}

	pekerjaan, err := u.pekerjaanRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if pekerjaan == nil {
		return errors.New("pekerjaan tidak ditemukan")
	}

	if pekerjaan.Status != entity.StatusAktif {
		return errors.New("hanya pekerjaan aktif yang dapat diselesaikan")
	}

	pekerjaan.Complete()
	return u.pekerjaanRepo.Update(ctx, id, pekerjaan)
}

func (u *PekerjaanAlumniUsecase) ResignJob(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("ID tidak valid")
	}

	pekerjaan, err := u.pekerjaanRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if pekerjaan == nil {
		return errors.New("pekerjaan tidak ditemukan")
	}

	if pekerjaan.Status != entity.StatusAktif {
		return errors.New("hanya pekerjaan aktif yang dapat diresign")
	}

	pekerjaan.Resign()
	return u.pekerjaanRepo.Update(ctx, id, pekerjaan)
}

func (u *PekerjaanAlumniUsecase) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("ID tidak valid")
	}

	// Check if pekerjaan exists
	existing, err := u.pekerjaanRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("pekerjaan tidak ditemukan")
	}

	return u.pekerjaanRepo.Delete(ctx, id)
}

func (u *PekerjaanAlumniUsecase) Search(ctx context.Context, query string, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return u.pekerjaanRepo.GetAll(ctx, limit, offset)
	}

	return u.pekerjaanRepo.Search(ctx, query, limit, offset)
}

func (u *PekerjaanAlumniUsecase) validatePekerjaan(pekerjaan *entity.PekerjaanAlumni) error {
	if pekerjaan.AlumniID == 0 {
		return errors.New("alumni ID harus diisi")
	}
	if pekerjaan.NamaCompany == "" {
		return errors.New("nama company harus diisi")
	}
	if pekerjaan.Posisi == "" {
		return errors.New("posisi harus diisi")
	}
	if pekerjaan.TanggalMulai.IsZero() {
		return errors.New("tanggal mulai harus diisi")
	}
	if pekerjaan.TanggalMulai.After(time.Now()) {
		return errors.New("tanggal mulai tidak boleh di masa depan")
	}
	return nil
}

func (u *PekerjaanAlumniUsecase) validatePekerjaanUpdate(pekerjaan *entity.PekerjaanAlumni) error {
	if pekerjaan.NamaCompany != "" && len(strings.TrimSpace(pekerjaan.NamaCompany)) == 0 {
		return errors.New("nama company tidak boleh kosong")
	}
	if pekerjaan.Posisi != "" && len(strings.TrimSpace(pekerjaan.Posisi)) == 0 {
		return errors.New("posisi tidak boleh kosong")
	}
	if !pekerjaan.TanggalMulai.IsZero() && pekerjaan.TanggalMulai.After(time.Now()) {
		return errors.New("tanggal mulai tidak boleh di masa depan")
	}
	return nil
}