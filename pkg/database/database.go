package database

import (
	"fmt"
	"log"

	"Fix-Go-Fiber-Backend/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabaseConnection(cfg *config.Config) (*gorm.DB, error) {
	var logLevel logger.LogLevel
	if cfg.App.Debug {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}

	var dialector gorm.Dialector
	
	switch cfg.Database.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.GetMySQLDSN())
	case "postgres":
		dialector = postgres.Open(cfg.GetPostgresDSN())
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
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

func RunMigrations(db *gorm.DB, cfg *config.Config) error {
	log.Println("Running database migrations...")
	
	// Use raw SQL migrations instead of GORM AutoMigrate
	err := CreateTables(db, cfg)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}