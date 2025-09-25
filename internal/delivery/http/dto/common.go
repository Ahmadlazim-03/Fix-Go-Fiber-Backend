package dto

import "time"

// Request DTOs
type CreateMahasiswaRequest struct {
	NIM      string `json:"nim" validate:"required,max=20"`
	Nama     string `json:"nama" validate:"required,max=100"`
	Jurusan  string `json:"jurusan" validate:"required,max=50"`
	Angkatan int    `json:"angkatan" validate:"required,min=1900,max=2100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateMahasiswaRequest struct {
	NIM      string `json:"nim,omitempty" validate:"omitempty,max=20"`
	Nama     string `json:"nama,omitempty" validate:"omitempty,max=100"`
	Jurusan  string `json:"jurusan,omitempty" validate:"omitempty,max=50"`
	Angkatan int    `json:"angkatan,omitempty" validate:"omitempty,min=1900,max=2100"`
	Email    string `json:"email,omitempty" validate:"omitempty,email,max=100"`
	Password string `json:"password,omitempty" validate:"omitempty,min=6"`
}

type CreateAlumniRequest struct {
	MahasiswaID uint   `json:"mahasiswa_id" validate:"required"`
	TahunLulus  int    `json:"tahun_lulus" validate:"required,min=1900,max=2100"`
	NoTelepon   string `json:"no_telepon,omitempty" validate:"omitempty,max=15"`
	Alamat      string `json:"alamat,omitempty"`
}

type UpdateAlumniRequest struct {
	MahasiswaID uint   `json:"mahasiswa_id,omitempty" validate:"omitempty"`
	TahunLulus  int    `json:"tahun_lulus,omitempty" validate:"omitempty,min=1900,max=2100"`
	NoTelepon   string `json:"no_telepon,omitempty" validate:"omitempty,max=15"`
	Alamat      string `json:"alamat,omitempty"`
}

type CreatePekerjaanAlumniRequest struct {
	AlumniID       uint      `json:"alumni_id" validate:"required"`
	NamaCompany    string    `json:"nama_company" validate:"required,max=100"`
	Posisi         string    `json:"posisi" validate:"required,max=100"`
	TanggalMulai   time.Time `json:"tanggal_mulai" validate:"required"`
	TanggalSelesai *time.Time `json:"tanggal_selesai,omitempty"`
	Status         string    `json:"status,omitempty" validate:"omitempty,oneof=aktif selesai resigned"`
	Deskripsi      string    `json:"deskripsi,omitempty"`
}

type UpdatePekerjaanAlumniRequest struct {
	AlumniID       uint       `json:"alumni_id,omitempty" validate:"omitempty"`
	NamaCompany    string     `json:"nama_company,omitempty" validate:"omitempty,max=100"`
	Posisi         string     `json:"posisi,omitempty" validate:"omitempty,max=100"`
	TanggalMulai   *time.Time `json:"tanggal_mulai,omitempty"`
	TanggalSelesai *time.Time `json:"tanggal_selesai,omitempty"`
	Status         string     `json:"status,omitempty" validate:"omitempty,oneof=aktif selesai resigned"`
	Deskripsi      string     `json:"deskripsi,omitempty"`
}

// Response DTOs
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Query parameters
type PaginationQuery struct {
	Page  int    `query:"page" validate:"omitempty,min=1"`
	Limit int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search string `query:"search" validate:"omitempty"`
}

func (p *PaginationQuery) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	return (p.Page - 1) * p.Limit
}

func (p *PaginationQuery) GetMeta(total int64) *Meta {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	
	totalPages := int(total) / p.Limit
	if int(total)%p.Limit > 0 {
		totalPages++
	}
	
	return &Meta{
		Page:       p.Page,
		Limit:      p.Limit,
		Total:      total,
		TotalPages: totalPages,
	}
}