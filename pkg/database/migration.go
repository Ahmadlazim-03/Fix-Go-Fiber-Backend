package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"Fix-Go-Fiber-Backend/pkg/config"

	"gorm.io/gorm"
)

// CreateTables creates all necessary tables using raw SQL for both PostgreSQL and MySQL
func CreateTables(db *gorm.DB, cfg *config.Config) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// First, drop existing tables if they have wrong structure
	if err := dropExistingTablesIfNeeded(sqlDB, cfg.Database.Driver); err != nil {
		log.Printf("Warning: failed to drop existing tables: %v", err)
	}

	var queries []string
	
	switch cfg.Database.Driver {
	case "postgres":
		queries = getPostgreSQLQueries()
	case "mysql":
		queries = getMySQLQueries()
	default:
		return fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	// Execute all queries
	for _, query := range queries {
		_, err := sqlDB.Exec(query)
		if err != nil {
			log.Printf("Error executing query: %s", query)
			return fmt.Errorf("failed to execute migration query: %w", err)
		}
	}

	log.Println("Database tables created successfully")
	
	// Create default admin user if not exists
	if err := createDefaultAdmin(sqlDB, cfg.Database.Driver); err != nil {
		log.Printf("Warning: failed to create default admin: %v", err)
	}

	return nil
}

func dropExistingTablesIfNeeded(sqlDB *sql.DB, driver string) error {
	// Check if tables exist and drop them to ensure clean migration
	var dropQueries []string
	
	switch driver {
	case "postgres":
		dropQueries = []string{
			`DROP TABLE IF EXISTS pekerjaan_alumni CASCADE`,
			`DROP TABLE IF EXISTS alumni CASCADE`,
			`DROP TABLE IF EXISTS admin_users CASCADE`,
			`DROP TABLE IF EXISTS mahasiswas CASCADE`,
		}
	case "mysql":
		dropQueries = []string{
			`DROP TABLE IF EXISTS pekerjaan_alumni`,
			`DROP TABLE IF EXISTS alumni`,
			`DROP TABLE IF EXISTS admin_users`,
			`DROP TABLE IF EXISTS mahasiswas`,
		}
	}
	
	for _, query := range dropQueries {
		_, err := sqlDB.Exec(query)
		if err != nil {
			log.Printf("Warning: failed to drop table: %v", err)
		}
	}
	
	log.Println("Existing tables dropped (if any)")
	return nil
}

func getPostgreSQLQueries() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS mahasiswas (
			id SERIAL PRIMARY KEY,
			nim VARCHAR(20) UNIQUE NOT NULL,
			nama VARCHAR(100) NOT NULL,
			jurusan VARCHAR(50) NOT NULL,
			angkatan INTEGER NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL
		)`,
		
		`CREATE TABLE IF NOT EXISTS alumni (
			id SERIAL PRIMARY KEY,
			mahasiswa_id INTEGER NOT NULL,
			tahun_lulus INTEGER NOT NULL,
			no_telepon VARCHAR(20),
			alamat TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL,
			FOREIGN KEY (mahasiswa_id) REFERENCES mahasiswas(id) ON DELETE CASCADE
		)`,
		
		`CREATE TABLE IF NOT EXISTS pekerjaan_alumni (
			id SERIAL PRIMARY KEY,
			alumni_id INTEGER NOT NULL,
			nama_company VARCHAR(100) NOT NULL,
			posisi VARCHAR(100) NOT NULL,
			tanggal_mulai DATE NOT NULL,
			tanggal_selesai DATE NULL,
			status VARCHAR(20) DEFAULT 'aktif',
			deskripsi TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL,
			FOREIGN KEY (alumni_id) REFERENCES alumni(id) ON DELETE CASCADE
		)`,
		
		`CREATE TABLE IF NOT EXISTS admin_users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(20) DEFAULT 'admin',
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL
		)`,

		// Create indexes for better performance
		`CREATE INDEX IF NOT EXISTS idx_mahasiswas_deleted_at ON mahasiswas(deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_mahasiswas_email ON mahasiswas(email)`,
		`CREATE INDEX IF NOT EXISTS idx_mahasiswas_nim ON mahasiswas(nim)`,
		`CREATE INDEX IF NOT EXISTS idx_alumni_deleted_at ON alumni(deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_alumni_mahasiswa_id ON alumni(mahasiswa_id)`,
		`CREATE INDEX IF NOT EXISTS idx_pekerjaan_alumni_deleted_at ON pekerjaan_alumni(deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_pekerjaan_alumni_alumni_id ON pekerjaan_alumni(alumni_id)`,
		`CREATE INDEX IF NOT EXISTS idx_admin_users_deleted_at ON admin_users(deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_admin_users_username ON admin_users(username)`,
		`CREATE INDEX IF NOT EXISTS idx_admin_users_email ON admin_users(email)`,
	}
}

func getMySQLQueries() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS mahasiswas (
			id INT AUTO_INCREMENT PRIMARY KEY,
			nim VARCHAR(20) UNIQUE NOT NULL,
			nama VARCHAR(100) NOT NULL,
			jurusan VARCHAR(50) NOT NULL,
			angkatan INT NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL
		)`,
		
		`CREATE TABLE IF NOT EXISTS alumni (
			id INT AUTO_INCREMENT PRIMARY KEY,
			mahasiswa_id INT NOT NULL,
			tahun_lulus INT NOT NULL,
			no_telepon VARCHAR(20),
			alamat TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL,
			FOREIGN KEY (mahasiswa_id) REFERENCES mahasiswas(id) ON DELETE CASCADE
		)`,
		
		`CREATE TABLE IF NOT EXISTS pekerjaan_alumni (
			id INT AUTO_INCREMENT PRIMARY KEY,
			alumni_id INT NOT NULL,
			nama_company VARCHAR(100) NOT NULL,
			posisi VARCHAR(100) NOT NULL,
			tanggal_mulai DATE NOT NULL,
			tanggal_selesai DATE NULL,
			status VARCHAR(20) DEFAULT 'aktif',
			deskripsi TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL,
			FOREIGN KEY (alumni_id) REFERENCES alumni(id) ON DELETE CASCADE
		)`,
		
		`CREATE TABLE IF NOT EXISTS admin_users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(20) DEFAULT 'admin',
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL
		)`,

		// Create indexes for better performance
		`CREATE INDEX IF NOT EXISTS idx_mahasiswas_deleted_at ON mahasiswas(deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_mahasiswas_email ON mahasiswas(email)`,
		`CREATE INDEX IF NOT EXISTS idx_mahasiswas_nim ON mahasiswas(nim)`,
		`CREATE INDEX IF NOT EXISTS idx_alumni_deleted_at ON alumni(deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_alumni_mahasiswa_id ON alumni(mahasiswa_id)`,
		`CREATE INDEX IF NOT EXISTS idx_pekerjaan_alumni_deleted_at ON pekerjaan_alumni(deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_pekerjaan_alumni_alumni_id ON pekerjaan_alumni(alumni_id)`,
		`CREATE INDEX IF NOT EXISTS idx_admin_users_deleted_at ON admin_users(deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_admin_users_username ON admin_users(username)`,
		`CREATE INDEX IF NOT EXISTS idx_admin_users_email ON admin_users(email)`,
	}
}

func createDefaultAdmin(sqlDB *sql.DB, driver string) error {
	// Check if admin already exists
	var count int
	var query string
	
	switch driver {
	case "postgres":
		query = `SELECT COUNT(*) FROM admin_users WHERE username = $1`
	case "mysql":
		query = `SELECT COUNT(*) FROM admin_users WHERE username = ?`
	}
	
	err := sqlDB.QueryRow(query, "admin").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing admin: %w", err)
	}
	
	if count > 0 {
		log.Println("Default admin already exists")
		return nil
	}

	// Create default admin (password: admin123)
	// Hash: $2a$12$DID7pK1MQljp0G3VQ1xUiekK1SXX1G04bYJy.bEN1zA6MxVaZ7eYC
	hashedPassword := "$2a$12$DID7pK1MQljp0G3VQ1xUiekK1SXX1G04bYJy.bEN1zA6MxVaZ7eYC"
	
	switch driver {
	case "postgres":
		query = `INSERT INTO admin_users (username, email, password, role, is_active, created_at, updated_at) 
				 VALUES ($1, $2, $3, $4, $5, $6, $7)`
	case "mysql":
		query = `INSERT INTO admin_users (username, email, password, role, is_active, created_at, updated_at) 
				 VALUES (?, ?, ?, ?, ?, ?, ?)`
	}
	
	now := time.Now()
	_, err = sqlDB.Exec(query, "admin", "admin@example.com", hashedPassword, "admin", true, now, now)
	if err != nil {
		return fmt.Errorf("failed to create default admin: %w", err)
	}
	
	log.Println("Default admin created successfully (username: admin, password: admin123)")
	return nil
}