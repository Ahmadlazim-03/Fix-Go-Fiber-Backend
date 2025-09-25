# Go Fiber Backend - Clean Architecture

Backend aplikasi untuk mengelola data mahasiswa, alumni, dan pekerjaan alumni menggunakan Go Fiber dengan implementasi Clean Architecture.

## 🏗️ Arsitektur

Aplikasi ini menggunakan **Clean Architecture** dengan struktur sebagai berikut:

```
# Fix-Go-Fiber-Backend

Role-based authentication and authorization system built with Go Fiber framework. This API provides secure endpoints for managing students (mahasiswa), alumni, and job records (pekerjaan) with JWT-based authentication and role-based access control.

## 🚀 Features

- **JWT Authentication** with role-based claims
- **Role-Based Access Control** (Admin, Alumni, Mahasiswa)
- **Soft Delete** functionality for all entities
- **Clean Architecture** implementation
- **PostgreSQL** database with GORM
- **Password Hashing** with bcrypt
- **Request Validation** with go-playground/validator
- **Structured Logging** with logrus
- **Auto Database Migration**

## 🏗️ Architecture

```
├── cmd/server/          # Application entry point
├── internal/            # Private application code
│   ├── entity/         # Domain models
│   ├── usecase/        # Business logic
│   ├── repository/     # Data access layer
│   └── delivery/       # HTTP handlers and routes
├── pkg/                # Reusable packages
│   ├── config/         # Configuration
│   ├── database/       # Database connection
│   ├── jwt/            # JWT utilities
│   ├── bcrypt/         # Password hashing
│   ├── validator/      # Request validation
│   └── logger/         # Logging
└── docs/               # Documentation
```

## 🛠️ Tech Stack

- **Go 1.21+**
- **Fiber v2** - Web framework
- **GORM** - ORM library
- **PostgreSQL** - Database
- **JWT** - Authentication tokens
- **Bcrypt** - Password hashing
- **Logrus** - Structured logging
- **Go Playground Validator** - Request validation

## ⚡ Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12+
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Ahmadlazim-03/Fix-Go-Fiber-Backend.git
   cd Fix-Go-Fiber-Backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` file:
   ```env
   APP_NAME=Fix-Go-Fiber-Backend
   APP_PORT=8080
   
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=fiber_db
   
   JWT_SECRET=your_super_secret_jwt_key_here
   LOG_LEVEL=info
   ```

4. **Setup PostgreSQL database**
   ```sql
   CREATE DATABASE fiber_db;
   ```

5. **Run database migrations**
   The application will automatically run migrations on startup, or you can use the reset script:
   ```bash
   # Windows PowerShell
   powershell -ExecutionPolicy Bypass -File .\scripts\reset_db.ps1
   
   # Manual SQL
   psql -U postgres -d fiber_db -f database/migrations/initial.sql
   ```

6. **Build and run the application**
   ```bash
   # Development
   go run cmd/server/main.go
   
   # Production build
   go build -o server cmd/server/main.go
   ./server  # or .\server.exe on Windows
   ```

The API will be available at `http://localhost:8080`

## 🔐 Authentication & Authorization

### User Roles

| Role | Description | Permissions |
|------|-------------|-------------|
| **Admin** | System administrator | Full CRUD access to all entities |
| **Alumni** | Graduated students | Manage own pekerjaan records |
| **Mahasiswa** | Current students | View and update own profile |

### JWT Token Structure

```json
{
  "user_id": "123",
  "role": "admin|alumni|mahasiswa",
  "exp": 1234567890
}
```

### Authentication Flow

1. **Login** with credentials → Receive JWT token
2. **Include token** in Authorization header: `Bearer <token>`
3. **Access protected endpoints** based on role permissions

## 📖 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Quick Test Endpoints

1. **Health Check**
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

2. **Register Mahasiswa**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/mahasiswa/register \
     -H "Content-Type: application/json" \
     -d '{"nim":"123456789","nama":"John Doe","email":"john@example.com","password":"password123","jurusan":"Teknik Informatika","angkatan":2020}'
   ```

3. **Admin Login**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/admin/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"admin123"}'
   ```

4. **Get Profile**
   ```bash
   curl -X GET http://localhost:8080/api/v1/auth/profile \
     -H "Authorization: Bearer <your_token>"
   ```

For complete API documentation, see [API_DOCUMENTATION.md](API_DOCUMENTATION.md)

## 🧪 Testing

### Using Postman

1. Import the Postman collection: `postman_collection.json`
2. Set up environment variables:
   - `base_url`: `http://localhost:8080/api/v1`
3. Run authentication requests to get tokens
4. Test protected endpoints with the tokens

### Using cURL

See the examples in [API_DOCUMENTATION.md](API_DOCUMENTATION.md#testing-the-api)

## 🗄️ Database Schema

### Entities

All entities include soft delete with `deleted_at` timestamp:

```sql
-- Mahasiswa (Students)
CREATE TABLE mahasiswa (
    id BIGSERIAL PRIMARY KEY,
    nim VARCHAR(20) UNIQUE NOT NULL,
    nama VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    jurusan VARCHAR(255) NOT NULL,
    angkatan INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Alumni (Graduates)
CREATE TABLE alumni (
    id BIGSERIAL PRIMARY KEY,
    nim VARCHAR(20) UNIQUE NOT NULL,
    nama VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    jurusan VARCHAR(255) NOT NULL,
    tahun_lulus INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Pekerjaan Alumni (Alumni Jobs)
CREATE TABLE pekerjaan_alumni (
    id BIGSERIAL PRIMARY KEY,
    alumni_id BIGINT REFERENCES alumni(id),
    nama_perusahaan VARCHAR(255) NOT NULL,
    posisi VARCHAR(255) NOT NULL,
    tahun_masuk INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Admin Users
CREATE TABLE admin_user (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | `Fix-Go-Fiber-Backend` |
| `APP_PORT` | Server port | `8080` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | `fiber_db` |
| `JWT_SECRET` | JWT signing secret | - |
| `LOG_LEVEL` | Log level | `info` |

### JWT Configuration

- **Algorithm**: HS256
- **Expiration**: Configurable (default: 24 hours)
- **Claims**: user_id, role, exp

### Logging Levels

- `trace` - Most detailed
- `debug` - Debug information
- `info` - General information
- `warn` - Warning messages
- `error` - Error messages
- `fatal` - Fatal errors

## 📁 Project Structure

```
Fix-Go-Fiber-Backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── entity/                  # Domain models
│   │   ├── mahasiswa.go
│   │   ├── alumni.go
│   │   ├── pekerjaan_alumni.go
│   │   └── admin_user.go
│   ├── usecase/                 # Business logic
│   │   ├── mahasiswa_usecase.go
│   │   ├── alumni_usecase.go
│   │   ├── pekerjaan_usecase.go
│   │   └── auth_service.go
│   ├── repository/              # Data access
│   │   └── postgres/
│   │       ├── mahasiswa_repository.go
│   │       ├── alumni_repository.go
│   │       ├── pekerjaan_repository.go
│   │       └── admin_repository.go
│   └── delivery/
│       └── http/
│           ├── handler/         # HTTP handlers
│           ├── route/           # Route definitions
│           └── middleware/      # HTTP middleware
├── pkg/                         # Reusable packages
│   ├── config/
│   ├── database/
│   ├── jwt/
│   ├── bcrypt/
│   ├── validator/
│   └── logger/
├── scripts/                     # Utility scripts
├── docs/                        # Documentation
├── .env.example                 # Environment template
├── go.mod                       # Go modules
├── go.sum                       # Dependencies
├── API_DOCUMENTATION.md         # Complete API docs
├── postman_collection.json      # Postman collection
└── README.md                    # This file
```

## 🚀 Deployment

### Docker (Coming Soon)

```dockerfile
# Dockerfile example
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/.env .
CMD ["./server"]
```

### Production Considerations

1. **Use environment-specific configurations**
2. **Implement proper logging aggregation**
3. **Set up database connection pooling**
4. **Use reverse proxy (nginx) for HTTPS**
5. **Implement rate limiting**
6. **Set up monitoring and health checks**

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

### Code Style

- Follow Go standard formatting (`go fmt`)
- Write tests for new features
- Update documentation
- Follow Clean Architecture principles

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Ahmad Lazim**
- GitHub: [@Ahmadlazim-03](https://github.com/Ahmadlazim-03)
- Email: your-email@example.com

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Web framework
- [GORM](https://gorm.io/) - ORM library
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation
- [Logrus](https://github.com/sirupsen/logrus) - Structured logging

## 📞 Support

If you have any questions or issues, please:

1. Check the [API Documentation](API_DOCUMENTATION.md)
2. Look at existing [Issues](https://github.com/Ahmadlazim-03/Fix-Go-Fiber-Backend/issues)
3. Create a new issue with detailed information

---

**Happy Coding! 🚀**
├── cmd/                        # Entry points
│   └── server/
│       └── main.go             # Main application
├── internal/                   # Application core
│   ├── domain/                 # Business entities & interfaces
│   │   ├── entity/             # Domain entities
│   │   ├── repository/         # Repository interfaces
│   │   └── service/            # Domain service interfaces
│   ├── usecase/                # Application business rules
│   ├── delivery/               # Presentation layer
│   │   └── http/
│   │       ├── handler/        # HTTP handlers
│   │       ├── middleware/     # HTTP middleware
│   │       ├── route/          # Route definitions
│   │       └── dto/            # Data transfer objects
│   └── repository/             # Infrastructure layer
│       ├── postgres/           # PostgreSQL implementations
│       ├── redis/              # Redis implementations
│       └── mock/               # Mock implementations
├── pkg/                        # Shared utilities
│   ├── config/                 # Configuration management
│   ├── database/               # Database connections
│   ├── logger/                 # Logging utilities
│   ├── jwt/                    # JWT utilities
│   ├── bcrypt/                 # Password hashing
│   └── utils/                  # General utilities
└── tests/                      # Test files
    ├── integration/            # Integration tests
    ├── unit/                   # Unit tests
    └── mocks/                  # Test mocks
```

## 🚀 Features

- ✅ Clean Architecture implementation
- ✅ RESTful API with Go Fiber
- ✅ PostgreSQL database integration
- ✅ JWT Authentication
- ✅ Password hashing with bcrypt
- ✅ Input validation
- ✅ Structured logging
- ✅ CORS support
- ✅ Auto migration
- ✅ Environment configuration

## 🛠️ Tech Stack

- **Framework**: Go Fiber v2
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT
- **Validation**: go-playground/validator
- **Logging**: Logrus
- **Configuration**: godotenv

## 📦 Installation

1. Clone the repository:
```bash
git clone https://github.com/Ahmadlazim-03/Fix-Go-Fiber-Backend.git
cd Fix-Go-Fiber-Backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Copy environment file:
```bash
cp .env.example .env
```

4. Configure your database and other settings in `.env`

5. Run the application:
```bash
go run cmd/server/main.go
```

Or build and run:
```bash
go build -o bin/server ./cmd/server
./bin/server
```

## 🔧 Configuration

Configure the application by editing the `.env` file:

```env
# Application Configuration
APP_NAME=Go-Fiber-Backend
APP_ENV=development
APP_HOST=localhost
APP_PORT=8080
APP_DEBUG=true

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=fiber_db
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Jakarta

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRE_MINUTES=60
```

## 📋 API Endpoints

### Health Check
- `GET /health` - Check application status

### Mahasiswa
- `POST /api/v1/mahasiswa` - Create new mahasiswa
- `GET /api/v1/mahasiswa` - Get all mahasiswa (with pagination)
- `GET /api/v1/mahasiswa/:id` - Get mahasiswa by ID
- `PUT /api/v1/mahasiswa/:id` - Update mahasiswa
- `DELETE /api/v1/mahasiswa/:id` - Delete mahasiswa

### Alumni (Coming Soon)
- Alumni management endpoints

### Pekerjaan Alumni (Coming Soon)
- Job management endpoints

## 🏛️ Database Schema

### Mahasiswa
- `id` (Primary Key)
- `nim` (Unique)
- `nama`
- `jurusan`
- `angkatan`
- `email` (Unique)
- `password` (Hashed)
- `created_at`, `updated_at`

### Alumni
- `id` (Primary Key)
- `mahasiswa_id` (Foreign Key)
- `tahun_lulus`
- `no_telepon`
- `alamat`
- `created_at`, `updated_at`

### Pekerjaan Alumni
- `id` (Primary Key)
- `alumni_id` (Foreign Key)
- `nama_company`
- `posisi`
- `tanggal_mulai`
- `tanggal_selesai`
- `status` (aktif/selesai/resigned)
- `deskripsi`
- `created_at`, `updated_at`

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/usecase
```

## 📁 Project Structure Principles

### Clean Architecture Layers:

1. **Domain Layer** (`internal/domain/`):
   - Entities: Core business objects
   - Repository Interfaces: Data access contracts
   - Service Interfaces: Domain service contracts

2. **Use Case Layer** (`internal/usecase/`):
   - Application business rules
   - Orchestrates data flow between entities

3. **Infrastructure Layer** (`internal/repository/`):
   - External interfaces implementations
   - Database access, external APIs

4. **Presentation Layer** (`internal/delivery/`):
   - HTTP handlers, middleware, routes
   - Input validation and response formatting

5. **Shared Layer** (`pkg/`):
   - Utilities, configurations, helpers
   - Can be imported by any layer

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.