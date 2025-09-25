package dto

// Pekerjaan Alumni DTOs (Updated for new unified mahasiswa structure)
type CreatePekerjaanRequest struct {
	// Reference mahasiswa (who must be alumni status)
	MahasiswaID *uint  `json:"mahasiswa_id" validate:"omitempty"`
	NIM         string `json:"nim" validate:"required_without=MahasiswaID,max=20"`
	
	NamaCompany    string `json:"nama_company" validate:"required,max=100"`
	Posisi         string `json:"posisi" validate:"required,max=100"`
	TanggalMulai   Date   `json:"tanggal_mulai" validate:"required"`
	TanggalSelesai *Date  `json:"tanggal_selesai" validate:"omitempty"`
	Status         string `json:"status" validate:"omitempty,oneof=aktif selesai resigned"`
	Deskripsi      string `json:"deskripsi" validate:"omitempty"`
}

type UpdatePekerjaanRequest struct {
	NamaCompany    string `json:"nama_company" validate:"omitempty,max=100"`
	Posisi         string `json:"posisi" validate:"omitempty,max=100"`
	TanggalMulai   *Date  `json:"tanggal_mulai" validate:"omitempty"`
	TanggalSelesai *Date  `json:"tanggal_selesai" validate:"omitempty"`
	Status         string `json:"status" validate:"omitempty,oneof=aktif selesai resigned"`
	Deskripsi      string `json:"deskripsi" validate:"omitempty"`
}

// Legacy Alumni DTOs - DEPRECATED
// Use MahasiswaService.Graduate() and MahasiswaService.UpdateAlumniData() instead