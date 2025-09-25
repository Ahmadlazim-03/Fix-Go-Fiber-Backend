package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/domain/repository"
	"Fix-Go-Fiber-Backend/internal/domain/dto"
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

// Implement service.PekerjaanAlumniService interface
func (u *PekerjaanAlumniUsecase) CreatePekerjaan(ctx context.Context, req *dto.CreatePekerjaanRequest) (*entity.PekerjaanAlumni, error) {
	// Check if alumni exists
	alumni, err := u.alumniRepo.GetByID(ctx, req.AlumniID)
	if err != nil {
		return nil, err
	}
	if alumni == nil {
		return nil, errors.New("alumni tidak ditemukan")
	}

	// Set default status if empty
	status := entity.StatusAktif
	if req.Status != "" {
		status = entity.StatusPekerjaan(req.Status)
	}

	pekerjaan := &entity.PekerjaanAlumni{
		AlumniID:       req.AlumniID,
		NamaCompany:    req.NamaCompany,
		Posisi:         req.Posisi,
		TanggalMulai:   req.TanggalMulai,
		TanggalSelesai: req.TanggalSelesai,
		Status:         status,
		Deskripsi:      req.Deskripsi,
	}

	if err := u.pekerjaanRepo.Create(ctx, pekerjaan); err != nil {
		return nil, err
	}

	return pekerjaan, nil
}

func (u *PekerjaanAlumniUsecase) GetPekerjaanByID(ctx context.Context, id uint) (*entity.PekerjaanAlumni, error) {
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

func (u *PekerjaanAlumniUsecase) GetPekerjaanByAlumniID(ctx context.Context, alumniID uint) ([]*entity.PekerjaanAlumni, error) {
	if alumniID == 0 {
		return nil, errors.New("alumni ID tidak valid")
	}

	return u.pekerjaanRepo.GetByAlumniID(ctx, alumniID)
}

func (u *PekerjaanAlumniUsecase) GetAllPekerjaan(ctx context.Context, search string, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	if search != "" {
		return u.pekerjaanRepo.Search(ctx, search, limit, offset)
	}

	return u.pekerjaanRepo.GetAll(ctx, limit, offset)
}

func (u *PekerjaanAlumniUsecase) UpdatePekerjaan(ctx context.Context, id uint, req *dto.UpdatePekerjaanRequest) (*entity.PekerjaanAlumni, error) {
	if id == 0 {
		return nil, errors.New("ID tidak valid")
	}

	// Check if pekerjaan exists
	existing, err := u.pekerjaanRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("pekerjaan tidak ditemukan")
	}

	// Update only non-empty fields
	if req.NamaCompany != "" {
		existing.NamaCompany = req.NamaCompany
	}
	if req.Posisi != "" {
		existing.Posisi = req.Posisi
	}
	if req.TanggalMulai != nil {
		existing.TanggalMulai = *req.TanggalMulai
	}
	if req.TanggalSelesai != nil {
		existing.TanggalSelesai = req.TanggalSelesai
	}
	if req.Status != "" {
		existing.Status = entity.StatusPekerjaan(req.Status)
	}
	if req.Deskripsi != "" {
		existing.Deskripsi = req.Deskripsi
	}

	if err := u.pekerjaanRepo.Update(ctx, id, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (u *PekerjaanAlumniUsecase) DeletePekerjaan(ctx context.Context, id uint) error {
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

	return u.pekerjaanRepo.Delete(ctx, id) // This will soft delete due to gorm.DeletedAt
}

// Legacy methods for backward compatibility
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