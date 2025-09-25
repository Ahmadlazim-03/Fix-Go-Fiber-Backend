package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Port        string
	Host        string
	Debug       bool
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	TimeZone string
}

type JWTConfig struct {
	SecretKey string
	Expire    string
}

type CORSConfig struct {
	AllowedOrigins     string
	AllowedMethods     string
	AllowedHeaders     string
	AllowCredentials   bool
}

func LoadConfig() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// .env file is optional
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	config := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "Go-Fiber-Backend"),
			Environment: getEnv("APP_ENV", "development"),
			Port:        getEnv("APP_PORT", "8080"),
			Host:        getEnv("APP_HOST", "localhost"),
			Debug:       getEnvAsBool("APP_DEBUG", true),
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "postgres"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "fiber_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			TimeZone: getEnv("DB_TIMEZONE", "Asia/Jakarta"),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET", "your-secret-key"),
			Expire:    getEnv("JWT_EXPIRE", "24h"),
		},
		CORS: CORSConfig{
			AllowedOrigins:   getEnv("CORS_ALLOWED_ORIGINS", "*"),
			AllowedMethods:   getEnv("CORS_ALLOWED_METHODS", "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS"),
			AllowedHeaders:   getEnv("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization"),
			AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", false),
		},
	}

	return config, nil
}

func (c *Config) GetPostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
		c.Database.TimeZone,
	)
}

func (c *Config) GetMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.TimeZone,
	)
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.App.Host, c.App.Port)
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultVal
}