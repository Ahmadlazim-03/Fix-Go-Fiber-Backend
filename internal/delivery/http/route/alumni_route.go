package route

import (
	"Fix-Go-Fiber-Backend/internal/delivery/http/handler"
	"Fix-Go-Fiber-Backend/internal/delivery/http/middleware"
	"Fix-Go-Fiber-Backend/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func SetupAlumniRoutes(
	api fiber.Router,
	alumniHandler *handler.AlumniHandler,
	jwtUtil *jwt.JWTUtil,
) {
	// Alumni routes
	alumni := api.Group("/alumni")
	
	// Admin only routes - Full CRUD access to all alumni data
	alumni.Post("/", middleware.RequireAuth(jwtUtil), middleware.AdminOnly(jwtUtil), alumniHandler.CreateAlumni)
	alumni.Get("/", middleware.RequireAuth(jwtUtil), middleware.AdminOnly(jwtUtil), alumniHandler.GetAllAlumni)
	alumni.Delete("/:id", middleware.RequireAuth(jwtUtil), middleware.AdminOnly(jwtUtil), alumniHandler.DeleteAlumni)
	
	// Admin or Self-access routes - Alumni can access their own data
	alumni.Get("/:id", middleware.RequireAuth(jwtUtil), middleware.AlumniOrAdmin(jwtUtil), alumniHandler.GetAlumniByID)
	alumni.Put("/:id", middleware.RequireAuth(jwtUtil), middleware.AlumniOrAdmin(jwtUtil), alumniHandler.UpdateAlumni)
}