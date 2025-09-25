package route

import (
	"Fix-Go-Fiber-Backend/internal/delivery/http/handler"
	"Fix-Go-Fiber-Backend/internal/delivery/http/middleware"
	"Fix-Go-Fiber-Backend/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	fiberMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRoutes(
	app *fiber.App,
	authHandler *handler.AuthHandler,
	mahasiswaHandler *handler.MahasiswaHandler,
	jwtUtil *jwt.JWTUtil,
) {
	// Global middleware
	app.Use(recover.New())
	app.Use(fiberMiddleware.New(middleware.NewLoggerMiddleware()))
	app.Use(middleware.NewCORSMiddleware())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Server is running",
		})
	})

	// API routes
	api := app.Group("/api/v1")
	
	// Auth routes (public)
	SetupAuthRoutes(api, authHandler, jwtUtil)
	
	// Protected routes
	SetupMahasiswaRoutes(api, mahasiswaHandler, jwtUtil)
	// TODO: Add when handlers are ready
	// SetupAlumniRoutes(api, alumniHandler, jwtUtil)
	// SetupPekerjaanAlumniRoutes(api, pekerjaanHandler, jwtUtil)
}