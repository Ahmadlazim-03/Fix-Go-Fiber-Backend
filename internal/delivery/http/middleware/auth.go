package middleware

import (
	"strings"

	"Fix-Go-Fiber-Backend/internal/domain/service"
	"Fix-Go-Fiber-Backend/pkg/jwt"
	"Fix-Go-Fiber-Backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// RoleBasedAuth creates a middleware that checks for specific roles
func RoleBasedAuth(jwtUtil *jwt.JWTUtil, allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Authorization header required"))
		}

		// Check if header starts with Bearer
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Invalid authorization format"))
		}

		token := tokenParts[1]
		
		// Validate token
		claims, err := jwtUtil.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Invalid or expired token"))
		}

		// Check if user role is allowed
		roleAllowed := false
		for _, allowedRole := range allowedRoles {
			if claims.Role == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			return c.Status(fiber.StatusForbidden).JSON(utils.ErrorResponse("Insufficient permissions"))
		}

		// Store user info in context
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("username", claims.Username)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// RequireAuth creates a middleware that just requires valid authentication
func RequireAuth(jwtUtil *jwt.JWTUtil) fiber.Handler {
	return RoleBasedAuth(jwtUtil, "mahasiswa", "alumni", "admin")
}

// AdminOnly creates a middleware that only allows admin access
func AdminOnly(jwtUtil *jwt.JWTUtil) fiber.Handler {
	return RoleBasedAuth(jwtUtil, "admin")
}

// AlumniOrAdmin creates a middleware that allows alumni or admin access
func AlumniOrAdmin(jwtUtil *jwt.JWTUtil) fiber.Handler {
	return RoleBasedAuth(jwtUtil, "alumni", "admin")
}

// GetUserFromContext extracts user claims from fiber context
func GetUserFromContext(c *fiber.Ctx) *service.JWTClaims {
	claims, ok := c.Locals("claims").(*service.JWTClaims)
	if !ok {
		return nil
	}
	return claims
}

// Legacy JWT middleware for backward compatibility
func JWTMiddleware(jwtManager *jwt.JWTUtil) fiber.Handler {
	return RequireAuth(jwtManager)
}