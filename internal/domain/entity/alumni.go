package entity

import (
	"time"
	"gorm.io/gorm"
)

type Alumni struct {
	ID             uint                `json:"id" gorm:"primaryKey"`
	MahasiswaID    uint                `json:"mahasiswa_id" gorm:"not null"`
	TahunLulus     int                 `json:"tahun_lulus" gorm:"not null"`
	NoTelepon      string              `json:"no_telepon" gorm:"size:15"`
	Alamat         string              `json:"alamat" gorm:"type:text"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	DeletedAt      gorm.DeletedAt      `json:"deleted_at" gorm:"index"`
	
	// Relations
	Mahasiswa      Mahasiswa           `json:"mahasiswa" gorm:"foreignKey:MahasiswaID"`
	PekerjaanList  []PekerjaanAlumni   `json:"pekerjaan_list,omitempty" gorm:"foreignKey:AlumniID"`
}

type AlumniResponse struct {
	ID          uint               `json:"id"`
	MahasiswaID uint               `json:"mahasiswa_id"`
	TahunLulus  int                `json:"tahun_lulus"`
	NoTelepon   string             `json:"no_telepon"`
	Alamat      string             `json:"alamat"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Mahasiswa   *MahasiswaResponse `json:"mahasiswa,omitempty"`
}

func (a *Alumni) ToResponse() *AlumniResponse {
	response := &AlumniResponse{
		ID:          a.ID,
		MahasiswaID: a.MahasiswaID,
		TahunLulus:  a.TahunLulus,
		NoTelepon:   a.NoTelepon,
		Alamat:      a.Alamat,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
	
	if a.Mahasiswa.ID != 0 {
		response.Mahasiswa = a.Mahasiswa.ToResponse()
	}
	
	return response
}

func (Alumni) TableName() string {
	return "alumni"
}