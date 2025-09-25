# API Documentation - Fix Go Fiber Backend

## 📋 Ringkasan Sistem

Sistem manajemen mahasiswa, alumni, dan pekerjaan dengan autentikasi berbasis peran (role-based authentication). Dibangun dengan Go Fiber dan PostgreSQL.

**Base URL**: `http://localhost:8080/api/v1`

---

## 👥 Jenis User & Hak Akses

| Role | Hak Akses |
|------|-----------|
| **Admin** | CRUD semua data (mahasiswa, alumni, pekerjaan) |
| **Alumni** | Kelola pekerjaan sendiri saja |
| **Mahasiswa** | Lihat/edit profil sendiri saja |

---

## 🔐 Cara Menggunakan API

### 1. **Daftar Akun Baru (Register)**

#### Daftar Mahasiswa
```bash
POST /auth/mahasiswa/register
```
```json
{
  "nim": "123456789",
  "nama": "John Doe",
  "email": "john@example.com", 
  "password": "password123",
  "jurusan": "Teknik Informatika",
  "angkatan": 2020
}
```

#### Daftar Alumni (Harus sudah jadi mahasiswa dulu)
```bash
POST /auth/alumni/register
```
```json
{
  "nim": "123456789",
  "nama": "John Doe",
  "email": "john@example.com",
  "password": "password123", 
  "jurusan": "Teknik Informatika",
  "tahun_lulus": 2022
}
```

### 2. **Login**

#### Login Mahasiswa
```bash
POST /auth/mahasiswa/login
```
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

#### Login Alumni  
```bash
POST /auth/alumni/login
```
```json
{
  "email": "john@example.com", 
  "password": "password123"
}
```

#### Login Admin
```bash
POST /auth/admin/login
```
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response Login:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": { /* data user */ },
    "role": "mahasiswa|alumni|admin",
    "expires_at": 1234567890
  }
}
```

### 3. **Akses API dengan Token**

Setelah login, gunakan token di header untuk akses API:
```bash
Authorization: Bearer <your_token_here>
```

---

## 📚 Endpoint API Lengkap

### 🔐 Authentication

| Method | Endpoint | Akses | Fungsi |
|--------|----------|-------|--------|
| POST | `/auth/mahasiswa/register` | Public | Daftar mahasiswa |
| POST | `/auth/alumni/register` | Public | Daftar alumni |
| POST | `/auth/mahasiswa/login` | Public | Login mahasiswa |
| POST | `/auth/alumni/login` | Public | Login alumni |
| POST | `/auth/admin/login` | Public | Login admin |
| GET | `/auth/profile` | Private | Lihat profil sendiri |

### 👨‍🎓 Mahasiswa

| Method | Endpoint | Akses | Fungsi |
|--------|----------|-------|--------|
| GET | `/mahasiswa` | Admin Only | Lihat semua mahasiswa |
| GET | `/mahasiswa/{id}` | Admin/Own | Lihat mahasiswa by ID |
| POST | `/mahasiswa` | Admin Only | Buat mahasiswa baru |
| PUT | `/mahasiswa/{id}` | Admin/Own | Update mahasiswa |
| DELETE | `/mahasiswa/{id}` | Admin Only | Hapus mahasiswa |

### 🎓 Alumni

| Method | Endpoint | Akses | Fungsi |
|--------|----------|-------|--------|
| GET | `/alumni` | Admin Only | Lihat semua alumni |
| GET | `/alumni/{id}` | Admin/Own | Lihat alumni by ID |
| POST | `/alumni` | Admin Only | Buat alumni baru |
| PUT | `/alumni/{id}` | Admin/Own | Update alumni |
| DELETE | `/alumni/{id}` | Admin Only | Hapus alumni |

### 💼 Pekerjaan Alumni

| Method | Endpoint | Akses | Fungsi |
|--------|----------|-------|--------|
| GET | `/pekerjaan` | Admin Only | Lihat semua pekerjaan |
| GET | `/alumni/{id}/pekerjaan` | Admin/Own | Pekerjaan by alumni |
| POST | `/pekerjaan` | Alumni/Admin | Buat pekerjaan baru |
| PUT | `/pekerjaan/{id}` | Alumni/Admin | Update pekerjaan |
| DELETE | `/pekerjaan/{id}` | Alumni/Admin | Hapus pekerjaan |

---

## 💼 Contoh Penggunaan Lengkap

### Step 1: Daftar Mahasiswa
```bash
curl -X POST http://localhost:8080/api/v1/auth/mahasiswa/register \
-H "Content-Type: application/json" \
-d '{
  "nim": "123456789",
  "nama": "John Doe", 
  "email": "john@example.com",
  "password": "password123",
  "jurusan": "Teknik Informatika",
  "angkatan": 2020
}'
```

### Step 2: Login Mahasiswa
```bash
curl -X POST http://localhost:8080/api/v1/auth/mahasiswa/login \
-H "Content-Type: application/json" \
-d '{
  "email": "john@example.com",
  "password": "password123" 
}'
```

### Step 3: Lihat Profil (Pakai Token)
```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
-H "Authorization: Bearer <token_dari_login>"
```

### Step 4: Daftar Alumni (Setelah Lulus)
```bash
curl -X POST http://localhost:8080/api/v1/auth/alumni/register \
-H "Content-Type: application/json" \
-d '{
  "nim": "123456789",
  "nama": "John Doe",
  "email": "john@example.com", 
  "password": "password123",
  "jurusan": "Teknik Informatika",
  "tahun_lulus": 2022
}'
```

### Step 5: Login Alumni & Buat Pekerjaan
```bash
# Login alumni
curl -X POST http://localhost:8080/api/v1/auth/alumni/login \
-H "Content-Type: application/json" \
-d '{
  "email": "john@example.com",
  "password": "password123"
}'

# Buat pekerjaan baru
curl -X POST http://localhost:8080/api/v1/pekerjaan \
-H "Authorization: Bearer <alumni_token>" \
-H "Content-Type: application/json" \
-d '{
  "nama_perusahaan": "Tech Corp",
  "posisi": "Software Engineer", 
  "tahun_masuk": 2023,
  "status": "aktif"
}'
```

---

## ⚠️ Error Responses

### 400 - Bad Request
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": [
    {
      "field": "email",
      "message": "email is required"
    }
  ]
}
```

### 401 - Unauthorized
```json
{
  "success": false,
  "message": "Invalid or missing token"
}
```

### 403 - Forbidden
```json
{
  "success": false,
  "message": "Insufficient permissions"
}
```

### 404 - Not Found
```json
{
  "success": false,
  "message": "Resource not found"
}
```

---

## 🎯 Flow Aplikasi

### 1. **Mahasiswa Baru:**
1. Daftar dengan `/auth/mahasiswa/register`
2. Login dengan `/auth/mahasiswa/login` 
3. Akses profil dengan `/auth/profile`
4. Update profil dengan `/mahasiswa/{id}`

### 2. **Setelah Lulus (Jadi Alumni):**
1. Daftar alumni dengan `/auth/alumni/register` (pakai NIM yang sama)
2. Login alumni dengan `/auth/alumni/login`
3. Buat pekerjaan dengan `/pekerjaan`
4. Kelola pekerjaan sendiri

### 3. **Admin:**
1. Login dengan `/auth/admin/login`
2. Kelola semua data mahasiswa, alumni, pekerjaan
3. CRUD operations pada semua entitas

---

## 🛠 Setup & Testing

### Environment Variables (.env)
```env
APP_NAME=Fix-Go-Fiber-Backend
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=fiber_db
JWT_SECRET=your_secret_key
```

### Quick Test
```bash
# Health check
curl http://localhost:8080/api/v1/health

# Daftar mahasiswa
curl -X POST http://localhost:8080/api/v1/auth/mahasiswa/register \
-H "Content-Type: application/json" \
-d '{"nim":"123456","nama":"Test User","email":"test@test.com","password":"123456","jurusan":"IT","angkatan":2023}'

# Login mahasiswa  
curl -X POST http://localhost:8080/api/v1/auth/mahasiswa/login \
-H "Content-Type: application/json" \
-d '{"email":"test@test.com","password":"123456"}'
```

---

## 📝 Database Schema

### mahasiswa
- id, nim (unique), nama, email (unique), password, jurusan, angkatan
- created_at, updated_at, deleted_at

### alumni  
- id, mahasiswa_id (FK), tahun_lulus, no_telepon, alamat
- created_at, updated_at, deleted_at

### pekerjaan_alumni
- id, alumni_id (FK), nama_perusahaan, posisi, tahun_masuk, status
- created_at, updated_at, deleted_at

### admin_user
- id, username (unique), email (unique), password
- created_at, updated_at, deleted_at

---

## ✨ Fitur Utama

- ✅ **Role-based Authentication** (Admin, Alumni, Mahasiswa)
- ✅ **JWT Tokens** dengan expiration
- ✅ **Public Registration** untuk mahasiswa & alumni  
- ✅ **Password Hashing** dengan bcrypt
- ✅ **Soft Delete** (data tidak hilang permanen)
- ✅ **Input Validation** otomatis
- ✅ **CORS Support** untuk frontend
- ✅ **Structured Logging** 
- ✅ **Auto Database Migration**

---

## 🎯 Kesimpulan

API ini memberikan sistem lengkap untuk:
- **Mahasiswa** bisa daftar, login, kelola profil
- **Alumni** bisa daftar, login, kelola pekerjaan  
- **Admin** bisa kelola semua data
- **Keamanan** dengan JWT dan role-based access
- **Data integrity** dengan soft delete

**Ready to use!** 🚀