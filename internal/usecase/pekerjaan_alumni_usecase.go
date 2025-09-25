package usecase

import (
	"context"
	"errors"
	"time"

	"Fix-Go-Fiber-Backend/internal/domain/dto"
	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/internal/domain/repository"
	"Fix-Go-Fiber-Backend/internal/domain/service"
)

type PekerjaanAlumniUsecase struct {
	pekerjaanRepo  repository.PekerjaanAlumniRepository
	mahasiswaRepo  repository.MahasiswaRepository
}

func NewPekerjaanAlumniUsecase(
	pekerjaanRepo repository.PekerjaanAlumniRepository,
	mahasiswaRepo repository.MahasiswaRepository,
) service.PekerjaanAlumniService {
	return &PekerjaanAlumniUsecase{
		pekerjaanRepo: pekerjaanRepo,
		mahasiswaRepo: mahasiswaRepo,
	}
}

// Implement service.PekerjaanAlumniService interface
func (u *PekerjaanAlumniUsecase) CreatePekerjaan(ctx context.Context, req *dto.CreatePekerjaanRequest) (*entity.PekerjaanAlumni, error) {
	var mahasiswaID uint
	
	// Jika mahasiswa_id disediakan, gunakan itu
	if req.MahasiswaID != nil && *req.MahasiswaID > 0 {
		mahasiswa, err := u.mahasiswaRepo.GetByID(ctx, *req.MahasiswaID)
		if err != nil {
			return nil, err
		}
		if mahasiswa == nil {
			return nil, errors.New("mahasiswa tidak ditemukan")
		}
		// Pastikan mahasiswa sudah alumni
		if !mahasiswa.IsAlumni() {
			return nil, errors.New("mahasiswa belum lulus, tidak dapat membuat data pekerjaan")
		}
		mahasiswaID = *req.MahasiswaID
	} else if req.NIM != "" {
		// Jika NIM disediakan, cari mahasiswa
		mahasiswa, err := u.mahasiswaRepo.GetByNIM(ctx, req.NIM)
		if err != nil || mahasiswa == nil {
			return nil, errors.New("mahasiswa dengan NIM tersebut tidak ditemukan")
		}
		
		// Pastikan mahasiswa sudah alumni
		if !mahasiswa.IsAlumni() {
			return nil, errors.New("mahasiswa belum lulus, tidak dapat membuat data pekerjaan")
		}
		
		mahasiswaID = mahasiswa.ID
	} else {
		return nil, errors.New("mahasiswa_id atau nim harus disediakan")
	}

	// Set default status if empty
	status := entity.StatusAktif
	if req.Status != "" {
		status = entity.StatusPekerjaan(req.Status)
	}

	// Convert Date types to time.Time
	var tanggalSelesai *time.Time
	if req.TanggalSelesai != nil {
		t := req.TanggalSelesai.Time
		tanggalSelesai = &t
	}

	pekerjaan := &entity.PekerjaanAlumni{
		MahasiswaID:    mahasiswaID,
		NamaCompany:    req.NamaCompany,
		Posisi:         req.Posisi,
		TanggalMulai:   req.TanggalMulai.Time,
		TanggalSelesai: tanggalSelesai,
		Status:         status,
		Deskripsi:      req.Deskripsi,
	}

	err := u.pekerjaanRepo.Create(ctx, pekerjaan)
	if err != nil {
		return nil, err
	}

	return pekerjaan, nil
}

func (u *PekerjaanAlumniUsecase) GetPekerjaanByID(ctx context.Context, id uint) (*entity.PekerjaanAlumni, error) {
	if id == 0 {
		return nil, errors.New("ID pekerjaan tidak valid")
	}

	return u.pekerjaanRepo.GetByID(ctx, id)
}

func (u *PekerjaanAlumniUsecase) GetPekerjaanByMahasiswaID(ctx context.Context, mahasiswaID uint) ([]*entity.PekerjaanAlumni, error) {
	if mahasiswaID == 0 {
		return nil, errors.New("ID mahasiswa tidak valid")
	}

	// Verifikasi mahasiswa exists dan alumni
	mahasiswa, err := u.mahasiswaRepo.GetByID(ctx, mahasiswaID)
	if err != nil {
		return nil, err
	}
	if mahasiswa == nil {
		return nil, errors.New("mahasiswa tidak ditemukan")
	}
	if !mahasiswa.IsAlumni() {
		return nil, errors.New("mahasiswa belum lulus, tidak ada data pekerjaan")
	}

	return u.pekerjaanRepo.GetByMahasiswaID(ctx, mahasiswaID)
}

func (u *PekerjaanAlumniUsecase) GetAllPekerjaan(ctx context.Context, search string, limit, offset int) ([]*entity.PekerjaanAlumni, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	if search != "" {
		// Implement search with filters
		filters := map[string]interface{}{
			"nama_company": search,
		}
		pekerjaans, err := u.pekerjaanRepo.GetWithFilters(ctx, filters)
		if err != nil {
			return nil, 0, err
		}
		return pekerjaans, int64(len(pekerjaans)), nil
	}

	return u.pekerjaanRepo.GetWithPagination(ctx, limit, offset)
}

func (u *PekerjaanAlumniUsecase) UpdatePekerjaan(ctx context.Context, id uint, req *dto.UpdatePekerjaanRequest) (*entity.PekerjaanAlumni, error) {
	if id == 0 {
		return nil, errors.New("ID pekerjaan tidak valid")
	}

	// Get existing pekerjaan
	existing, err := u.pekerjaanRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("pekerjaan tidak ditemukan")
	}

	// Update fields if provided
	if req.NamaCompany != "" {
		existing.NamaCompany = req.NamaCompany
	}
	if req.Posisi != "" {
		existing.Posisi = req.Posisi
	}
	if req.TanggalMulai != nil {
		existing.TanggalMulai = req.TanggalMulai.Time
	}
	if req.TanggalSelesai != nil {
		t := req.TanggalSelesai.Time
		existing.TanggalSelesai = &t
	}
	if req.Status != "" {
		existing.Status = entity.StatusPekerjaan(req.Status)
	}
	if req.Deskripsi != "" {
		existing.Deskripsi = req.Deskripsi
	}

	err = u.pekerjaanRepo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (u *PekerjaanAlumniUsecase) DeletePekerjaan(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("ID pekerjaan tidak valid")
	}

	// Check if exists
	existing, err := u.pekerjaanRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("pekerjaan tidak ditemukan")
	}

	return u.pekerjaanRepo.Delete(ctx, id)
}