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
	AlumniID     uint            `json:"alumni_id" gorm:"not null"`
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
	Alumni Alumni `json:"alumni,omitempty" gorm:"foreignKey:AlumniID"`
}

type PekerjaanAlumniResponse struct {
	ID             uint               `json:"id"`
	AlumniID       uint               `json:"alumni_id"`
	NamaCompany    string             `json:"nama_company"`
	Posisi         string             `json:"posisi"`
	TanggalMulai   time.Time          `json:"tanggal_mulai"`
	TanggalSelesai *time.Time         `json:"tanggal_selesai"`
	Status         StatusPekerjaan    `json:"status"`
	Deskripsi      string             `json:"deskripsi"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	Alumni         *AlumniResponse    `json:"alumni,omitempty"`
}

func (p *PekerjaanAlumni) ToResponse() *PekerjaanAlumniResponse {
	response := &PekerjaanAlumniResponse{
		ID:             p.ID,
		AlumniID:       p.AlumniID,
		NamaCompany:    p.NamaCompany,
		Posisi:         p.Posisi,
		TanggalMulai:   p.TanggalMulai,
		TanggalSelesai: p.TanggalSelesai,
		Status:         p.Status,
		Deskripsi:      p.Deskripsi,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
	
	if p.Alumni.ID != 0 {
		response.Alumni = p.Alumni.ToResponse()
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