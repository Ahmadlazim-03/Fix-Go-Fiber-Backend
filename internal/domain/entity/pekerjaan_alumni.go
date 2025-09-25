package entity

import (
	"time"
	"gorm.io/gorm"
)

type StatusPekerjaan string

const (
	StatusAktif    StatusPekerjaan = "aktif"
	StatusSelesai  StatusPekerjaan = "selesai"
	StatusResigned StatusPekerjaan = "resigned"
)

type PekerjaanAlumni struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	MahasiswaID  uint            `json:"mahasiswa_id" gorm:"not null"` // Reference to Mahasiswa (who is alumni)
	NamaCompany  string          `json:"nama_company" gorm:"not null;size:100"`
	Posisi       string          `json:"posisi" gorm:"not null;size:100"`
	TanggalMulai time.Time       `json:"tanggal_mulai" gorm:"not null"`
	TanggalSelesai *time.Time    `json:"tanggal_selesai"`
	Status       StatusPekerjaan `json:"status" gorm:"type:varchar(20);default:'aktif'"`
	Deskripsi    string          `json:"deskripsi" gorm:"type:text"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `json:"deleted_at" gorm:"index"`
	
	// Relations
	Mahasiswa Mahasiswa `json:"mahasiswa,omitempty" gorm:"foreignKey:MahasiswaID"`
}

type PekerjaanAlumniResponse struct {
	ID             uint               `json:"id"`
	MahasiswaID    uint               `json:"mahasiswa_id"`
	NamaCompany    string             `json:"nama_company"`
	Posisi         string             `json:"posisi"`
	TanggalMulai   time.Time          `json:"tanggal_mulai"`
	TanggalSelesai *time.Time         `json:"tanggal_selesai"`
	Status         StatusPekerjaan    `json:"status"`
	Deskripsi      string             `json:"deskripsi"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	Mahasiswa      *MahasiswaResponse `json:"mahasiswa,omitempty"`
}

func (p *PekerjaanAlumni) ToResponse() *PekerjaanAlumniResponse {
	response := &PekerjaanAlumniResponse{
		ID:             p.ID,
		MahasiswaID:    p.MahasiswaID,
		NamaCompany:    p.NamaCompany,
		Posisi:         p.Posisi,
		TanggalMulai:   p.TanggalMulai,
		TanggalSelesai: p.TanggalSelesai,
		Status:         p.Status,
		Deskripsi:      p.Deskripsi,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
	
	if p.Mahasiswa.ID != 0 {
		response.Mahasiswa = p.Mahasiswa.ToResponse()
	}
	
	return response
}

func (p *PekerjaanAlumni) IsActive() bool {
	return p.Status == StatusAktif
}

func (p *PekerjaanAlumni) Complete() {
	p.Status = StatusSelesai
	now := time.Now()
	p.TanggalSelesai = &now
}

func (p *PekerjaanAlumni) Resign() {
	p.Status = StatusResigned
	now := time.Now()
	p.TanggalSelesai = &now
}

func (PekerjaanAlumni) TableName() string {
	return "pekerjaan_alumni"
}