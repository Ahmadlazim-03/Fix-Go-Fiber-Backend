package entity

import (
	"time"
	"gorm.io/gorm"
)

// Status mahasiswa
type StatusMahasiswa string

const (
	StatusMahasiswaActive     StatusMahasiswa = "active"      // Masih kuliah
	StatusMahasiswaGraduated  StatusMahasiswa = "graduated"   // Sudah lulus (alumni)
	StatusMahasiswaDroppedOut StatusMahasiswa = "dropped_out" // Drop out
	StatusMahasiswaSuspended  StatusMahasiswa = "suspended"   // Diskors
)

// Unified Mahasiswa entity - covers both mahasiswa and alumni
type Mahasiswa struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	NIM       string         `json:"nim" gorm:"unique;not null;size:20"`
	Nama      string         `json:"nama" gorm:"not null;size:100"`
	Jurusan   string         `json:"jurusan" gorm:"not null;size:50"`
	Angkatan  int            `json:"angkatan" gorm:"not null"`
	Email     string         `json:"email" gorm:"unique;not null;size:100"`
	Password  string         `json:"-" gorm:"not null"`
	
	// Status Evolution
	Status    StatusMahasiswa `json:"status" gorm:"type:varchar(20);default:'active'"`
	
	// Alumni Fields (optional - filled when graduated)
	TahunLulus    *int    `json:"tahun_lulus" gorm:"null"` // NULL jika belum lulus
	NoTelepon     string  `json:"no_telepon" gorm:"size:15"`
	AlamatAlumni  string  `json:"alamat_alumni" gorm:"type:text"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	
	// Relations
	PekerjaanList []PekerjaanAlumni `json:"pekerjaan_list,omitempty" gorm:"foreignKey:MahasiswaID"`
}

// Helper methods
func (m *Mahasiswa) IsAlumni() bool {
	return m.Status == StatusMahasiswaGraduated
}

func (m *Mahasiswa) IsActive() bool {
	return m.Status == StatusMahasiswaActive
}

// Graduate - convert mahasiswa to alumni status
func (m *Mahasiswa) Graduate(tahunLulus int, noTelepon, alamat string) {
	m.Status = StatusMahasiswaGraduated
	m.TahunLulus = &tahunLulus
	m.NoTelepon = noTelepon
	m.AlamatAlumni = alamat
}

type MahasiswaResponse struct {
	ID            uint            `json:"id"`
	NIM           string          `json:"nim"`
	Nama          string          `json:"nama"`
	Jurusan       string          `json:"jurusan"`
	Angkatan      int             `json:"angkatan"`
	Email         string          `json:"email"`
	Status        StatusMahasiswa `json:"status"`
	TahunLulus    *int            `json:"tahun_lulus,omitempty"`
	NoTelepon     string          `json:"no_telepon,omitempty"`
	AlamatAlumni  string          `json:"alamat_alumni,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (m *Mahasiswa) ToResponse() *MahasiswaResponse {
	return &MahasiswaResponse{
		ID:            m.ID,
		NIM:           m.NIM,
		Nama:          m.Nama,
		Jurusan:       m.Jurusan,
		Angkatan:      m.Angkatan,
		Email:         m.Email,
		Status:        m.Status,
		TahunLulus:    m.TahunLulus,
		NoTelepon:     m.NoTelepon,
		AlamatAlumni:  m.AlamatAlumni,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func (Mahasiswa) TableName() string {
	return "mahasiswa"
}

// Model Alumni (legacy/optional)
type Alumni struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	MahasiswaID  uint      `json:"mahasiswa_id" gorm:"not null;index"`
	TahunLulus   int       `json:"tahun_lulus" gorm:"not null"`
	NoTelepon    string    `json:"no_telepon" gorm:"size:15"`
	AlamatAlumni string    `json:"alamat_alumni" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Alumni) TableName() string {
	return "alumni"
}