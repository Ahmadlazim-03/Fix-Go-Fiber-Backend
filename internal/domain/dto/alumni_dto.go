package dto

import "time"

// Alumni DTOs
type CreateAlumniRequest struct {
	MahasiswaID uint   `json:"mahasiswa_id" validate:"required"`
	TahunLulus  int    `json:"tahun_lulus" validate:"required,min=1900,max=2100"`
	NoTelepon   string `json:"no_telepon" validate:"omitempty,max=15"`
	Alamat      string `json:"alamat" validate:"omitempty"`
}

type UpdateAlumniRequest struct {
	TahunLulus int    `json:"tahun_lulus" validate:"omitempty,min=1900,max=2100"`
	NoTelepon  string `json:"no_telepon" validate:"omitempty,max=15"`
	Alamat     string `json:"alamat" validate:"omitempty"`
}

// Pekerjaan Alumni DTOs
type CreatePekerjaanRequest struct {
	AlumniID       uint      `json:"alumni_id" validate:"required"`
	NamaCompany    string    `json:"nama_company" validate:"required,max=100"`
	Posisi         string    `json:"posisi" validate:"required,max=100"`
	TanggalMulai   time.Time `json:"tanggal_mulai" validate:"required"`
	TanggalSelesai *time.Time `json:"tanggal_selesai" validate:"omitempty"`
	Status         string    `json:"status" validate:"omitempty,oneof=aktif selesai resigned"`
	Deskripsi      string    `json:"deskripsi" validate:"omitempty"`
}

type UpdatePekerjaanRequest struct {
	NamaCompany    string     `json:"nama_company" validate:"omitempty,max=100"`
	Posisi         string     `json:"posisi" validate:"omitempty,max=100"`
	TanggalMulai   *time.Time `json:"tanggal_mulai" validate:"omitempty"`
	TanggalSelesai *time.Time `json:"tanggal_selesai" validate:"omitempty"`
	Status         string     `json:"status" validate:"omitempty,oneof=aktif selesai resigned"`
	Deskripsi      string     `json:"deskripsi" validate:"omitempty"`
}