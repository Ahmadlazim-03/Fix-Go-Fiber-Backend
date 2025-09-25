package route

import (
	"Fix-Go-Fiber-Backend/internal/delivery/http/handler"
	"Fix-Go-Fiber-Backend/internal/delivery/http/middleware"
	"Fix-Go-Fiber-Backend/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app fiber.Router, authHandler *handler.AuthHandler, jwtUtil *jwt.JWTUtil) {
	auth := app.Group("/auth")

	// Public auth routes - Registration
	auth.Post("/mahasiswa/register", authHandler.RegisterMahasiswa)
	auth.Post("/alumni/register", authHandler.RegisterAlumni)

	// Public auth routes - Login
	auth.Post("/mahasiswa/login", authHandler.LoginMahasiswa)
	auth.Post("/alumni/login", authHandler.LoginAlumni)
	auth.Post("/admin/login", authHandler.LoginAdmin)

	// Protected profile route
	auth.Get("/profile", middleware.RequireAuth(jwtUtil), authHandler.GetProfile)
}