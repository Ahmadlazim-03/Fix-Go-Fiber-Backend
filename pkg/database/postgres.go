package database

import (
	"fmt"
	"log"

	"Fix-Go-Fiber-Backend/internal/domain/entity"
	"Fix-Go-Fiber-Backend/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresConnection(cfg *config.Config) (*gorm.DB, error) {
	var logLevel logger.LogLevel
	if cfg.App.Debug {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}

	db, err := gorm.Open(postgres.Open(cfg.GetDatabaseDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")
	
	err := db.AutoMigrate(
		&entity.Mahasiswa{},
		&entity.Alumni{},
		&entity.PekerjaanAlumni{},
		&entity.AdminUser{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}