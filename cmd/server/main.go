package main

import (
	"log"

	"Fix-Go-Fiber-Backend/internal/delivery/http/handler"
	"Fix-Go-Fiber-Backend/internal/delivery/http/route"
	"Fix-Go-Fiber-Backend/internal/repository/postgres"
	"Fix-Go-Fiber-Backend/internal/usecase"
	"Fix-Go-Fiber-Backend/pkg/bcrypt"
	"Fix-Go-Fiber-Backend/pkg/config"
	"Fix-Go-Fiber-Backend/pkg/database"
	"Fix-Go-Fiber-Backend/pkg/jwt"
	"Fix-Go-Fiber-Backend/pkg/logger"
	"Fix-Go-Fiber-Backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Setup logger
	appLogger := logger.NewLogrus(cfg)
	appLogger.Info("Starting application...")

	// Connect to database
	db, err := database.NewDatabaseConnection(cfg)
	if err != nil {
		appLogger.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.RunMigrations(db, cfg); err != nil {
		appLogger.Fatal("Failed to run migrations:", err)
	}

	// Initialize utilities
	bcryptHelper := bcrypt.NewBcryptHelper(12)
	bcryptUtil := bcrypt.NewBcryptUtil(12)
	jwtUtil := jwt.NewJWTUtil(cfg)
	customValidator := validator.NewCustomValidator()
	standardValidator := customValidator.GetValidator() // Get standard validator for mahasiswa handler

	// Initialize repositories
	mahasiswaRepo := postgres.NewMahasiswaRepository(db)
	alumniRepo := postgres.NewAlumniRepository(db)
	adminRepo := postgres.NewAdminUserRepository(db)
	pekerjaanAlumniRepo := postgres.NewPekerjaanAlumniRepository(db)

	// Initialize use cases
	mahasiswaUsecase := usecase.NewMahasiswaUsecase(mahasiswaRepo, bcryptHelper)
	alumniUsecase := usecase.NewAlumniUsecase(alumniRepo, mahasiswaRepo)
	pekerjaanUsecase := usecase.NewPekerjaanAlumniUsecase(pekerjaanAlumniRepo, alumniRepo)
	authService := usecase.NewAuthService(mahasiswaRepo, alumniRepo, adminRepo, jwtUtil, bcryptUtil)

	// Initialize handlers
	mahasiswaHandler := handler.NewMahasiswaHandler(mahasiswaUsecase, standardValidator)
	alumniHandler := handler.NewAlumniHandler(alumniUsecase, standardValidator)
	pekerjaanHandler := handler.NewPekerjaanAlumniHandler(pekerjaanUsecase, standardValidator)
	authHandler := handler.NewAuthHandler(authService, customValidator)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: cfg.App.Name,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			appLogger.Error("Request failed:", err)
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	// Setup routes
	route.SetupRoutes(app, cfg, authHandler, mahasiswaHandler, alumniHandler, pekerjaanHandler, jwtUtil)

	// Start server
	address := ":" + cfg.App.Port
	appLogger.Infof("Server starting on %s", address)
	if err := app.Listen(address); err != nil {
		appLogger.Fatal("Failed to start server:", err)
	}
}