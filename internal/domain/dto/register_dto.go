package dto

// Register DTOs
type RegisterMahasiswaRequest struct {
	NIM      string `json:"nim" validate:"required"`
	Nama     string `json:"nama" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Jurusan  string `json:"jurusan" validate:"required"`
	Angkatan int    `json:"angkatan" validate:"required,min=1900,max=2030"`
}

type RegisterAlumniRequest struct {
	NIM        string `json:"nim" validate:"required"`
	Nama       string `json:"nama" validate:"required,min=2"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
	Jurusan    string `json:"jurusan" validate:"required"`
	TahunLulus int    `json:"tahun_lulus" validate:"required,min=1900,max=2030"`
}

type RegisterResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}