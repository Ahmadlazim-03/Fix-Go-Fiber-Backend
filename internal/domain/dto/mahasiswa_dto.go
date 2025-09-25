package dto

// Mahasiswa DTOs
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