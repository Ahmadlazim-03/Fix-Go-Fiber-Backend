package route

import (
	"Fix-Go-Fiber-Backend/internal/delivery/http/handler"
	"Fix-Go-Fiber-Backend/internal/delivery/http/middleware"
	"Fix-Go-Fiber-Backend/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func SetupMahasiswaRoutes(app fiber.Router, handler *handler.MahasiswaHandler, jwtUtil *jwt.JWTUtil) {
	mahasiswa := app.Group("/mahasiswa")

	// Public routes
	mahasiswa.Post("/", handler.Create)

	// Admin only routes
	mahasiswa.Get("/", middleware.AdminOnly(jwtUtil), handler.GetAll)
	mahasiswa.Delete("/:id", middleware.AdminOnly(jwtUtil), handler.Delete)

	// Admin or own record routes (mahasiswa can view/update their own record)
	mahasiswa.Get("/:id", middleware.RoleBasedAuth(jwtUtil, "mahasiswa", "alumni", "admin"), handler.GetByID)
	mahasiswa.Put("/:id", middleware.RoleBasedAuth(jwtUtil, "mahasiswa", "alumni", "admin"), handler.Update)
}