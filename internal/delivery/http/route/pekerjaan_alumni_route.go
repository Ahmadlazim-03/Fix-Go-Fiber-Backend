package route

import (
	"Fix-Go-Fiber-Backend/internal/delivery/http/handler"
	"Fix-Go-Fiber-Backend/internal/delivery/http/middleware"
	"Fix-Go-Fiber-Backend/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func SetupPekerjaanAlumniRoutes(
	api fiber.Router,
	pekerjaanHandler *handler.PekerjaanAlumniHandler,
	jwtUtil *jwt.JWTUtil,
) {
	// Pekerjaan Alumni routes
	pekerjaan := api.Group("/pekerjaan")
	
	// Admin only routes - Full CRUD access to all pekerjaan data
	pekerjaan.Get("/", middleware.RequireAuth(jwtUtil), middleware.AdminOnly(jwtUtil), pekerjaanHandler.GetAllPekerjaan)
	
	// Alumni and Admin routes - Alumni can manage their own pekerjaan, Admin can manage any
	pekerjaan.Post("/", middleware.RequireAuth(jwtUtil), middleware.AlumniOrAdmin(jwtUtil), pekerjaanHandler.CreatePekerjaan)
	pekerjaan.Get("/:id", middleware.RequireAuth(jwtUtil), middleware.AlumniOrAdmin(jwtUtil), pekerjaanHandler.GetPekerjaanByID)
	pekerjaan.Put("/:id", middleware.RequireAuth(jwtUtil), middleware.AlumniOrAdmin(jwtUtil), pekerjaanHandler.UpdatePekerjaan)
	pekerjaan.Delete("/:id", middleware.RequireAuth(jwtUtil), middleware.AlumniOrAdmin(jwtUtil), pekerjaanHandler.DeletePekerjaan)
	
	// Get pekerjaan by alumni ID - Alumni can get their own, Admin can get any
	pekerjaan.Get("/alumni/:alumni_id", middleware.RequireAuth(jwtUtil), middleware.AlumniOrAdmin(jwtUtil), pekerjaanHandler.GetPekerjaanByAlumniID)
}