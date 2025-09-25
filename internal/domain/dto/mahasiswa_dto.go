package dto

// Mahasiswa registration (initial)
type CreateMahasiswaRequest struct {
	NIM      string `json:"nim" validate:"required,max=20"`
	Nama     string `json:"nama" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6"`
	Jurusan  string `json:"jurusan" validate:"required,max=50"`
	Angkatan int    `json:"angkatan" validate:"required,min=1900,max=2100"`
}

// Update mahasiswa profile (while active)
type UpdateMahasiswaRequest struct {
	Nama      string `json:"nama,omitempty" validate:"omitempty,max=100"`
	Email     string `json:"email,omitempty" validate:"omitempty,email,max=100"`
	NoTelepon string `json:"no_telepon,omitempty" validate:"omitempty,max=15"`
}

// Graduate mahasiswa to alumni status
type GraduateMahasiswaRequest struct {
	MahasiswaID   uint   `json:"mahasiswa_id" validate:"required"`
	TahunLulus    int    `json:"tahun_lulus" validate:"required,min=1900,max=2100"`
	NoTelepon     string `json:"no_telepon" validate:"omitempty,max=15"`
	AlamatAlumni  string `json:"alamat_alumni" validate:"omitempty"`
}

// Update alumni data (after graduation)
type UpdateAlumniDataRequest struct {
	NoTelepon    string `json:"no_telepon,omitempty" validate:"omitempty,max=15"`
	AlamatAlumni string `json:"alamat_alumni,omitempty" validate:"omitempty"`
}

// Query filters
type MahasiswaListRequest struct {
	Status   string `query:"status"`   // active, graduated, dropped_out
	Jurusan  string `query:"jurusan"`  // filter by jurusan
	Angkatan *int   `query:"angkatan"` // filter by angkatan
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
}

// Alumni list (just graduated mahasiswa)
type AlumniListRequest struct {
	Jurusan    string `query:"jurusan"`     // filter by jurusan
	TahunLulus *int   `query:"tahun_lulus"` // filter by graduation year
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
}