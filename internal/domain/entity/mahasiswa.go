package entity

import (
	"time"
	"gorm.io/gorm"
)

type Mahasiswa struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	NIM       string         `json:"nim" gorm:"unique;not null;size:20"`
	Nama      string         `json:"nama" gorm:"not null;size:100"`
	Jurusan   string         `json:"jurusan" gorm:"not null;size:50"`
	Angkatan  int            `json:"angkatan" gorm:"not null"`
	Email     string         `json:"email" gorm:"unique;not null;size:100"`
	Password  string         `json:"-" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type MahasiswaResponse struct {
	ID        uint      `json:"id"`
	NIM       string    `json:"nim"`
	Nama      string    `json:"nama"`
	Jurusan   string    `json:"jurusan"`
	Angkatan  int       `json:"angkatan"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Mahasiswa) ToResponse() *MahasiswaResponse {
	return &MahasiswaResponse{
		ID:        m.ID,
		NIM:       m.NIM,
		Nama:      m.Nama,
		Jurusan:   m.Jurusan,
		Angkatan:  m.Angkatan,
		Email:     m.Email,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (Mahasiswa) TableName() string {
	return "mahasiswa"
}